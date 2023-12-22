//  Copyright (C) 晓白齐齐,版权所有. 2023

package hcore

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"unsafe"

	"github.com/bqqsrc/goper/http"
)

// stringToBytes converts string to byte slice without a memory allocation.
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func stringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// bytesToString converts byte slice to string without a memory allocation.
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func bytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

var (
	strColon = []byte(":")
	strStar  = []byte("*")
	strSlash = []byte("/")
)

func countParams(path string) uint16 {
	var n uint16
	s := stringToBytes(path)
	n += uint16(bytes.Count(s, strColon))
	n += uint16(bytes.Count(s, strStar))
	return n
}

func countSections(path string) uint16 {
	s := stringToBytes(path)
	return uint16(bytes.Count(s, strSlash))
}

type nodeType = uint8

const domainSep uint8 = '.' // 域名的分割字符
const maxSection = 10       // 最大的节点数
const maxParamsCount = 10   // 最大的参数数量

const (
	static   nodeType = iota // 静态节点
	root                     // 根节点
	param                    // 参数节点
	catchAll                 // 全匹配节点
)

type node struct {
	path      string           // 节点的路径
	indices   string           // 本节点的所有静态子节点路径的首字符，下标索引与children下标索引对应
	nType     nodeType         // 节点类型，默认为static
	priority  uint32           // 以本节点为根节点的树所配置的路由数量
	children  []*node          // 子节点
	wildChild bool             // 是否有参数子节点或者全匹配子节点
	fullPath  string           // 全路径
	handler   http.HttpHandler // 存储的目标数据 Value any
}

// 添加路由
func (n *node) addRouter(path string, handler http.HttpHandler) error {
	if path == "" {
		return fmt.Errorf("path must not be empty")
	}
	if path[0] != '/' && path[0] != domainSep {
		return fmt.Errorf("a router path must begin with '/' or '%s', path: %s", string(domainSep), path)
	}
	fullPath := path
	n.priority++

	// 空节点
	if len(n.path) == 0 && len(n.children) == 0 {
		if err := n.insertChild(path, fullPath, handler); err != nil {
			return err
		}
		n.nType = root
		return nil
	}

	parentFullPathIndex := 0
walk:
	for {
		// 查找路径与当前节点路径最长的一致前缀
		// 包括 ： 和 * 字符
		i := longestCommonPrefix(path, n.path)

		// 最长前缀比当前节点路径要短
		// 需要将当前节点进行拆分，以最长前缀和后缀两部分拆分成两个节点
		if i < len(n.path) {
			// 新的子节点继承大部分原节点的数据
			child := node{
				path:      n.path[i:],
				wildChild: n.wildChild,
				nType:     static,
				indices:   n.indices,
				children:  n.children,
				handler:   n.handler,
				priority:  n.priority - 1,
				fullPath:  n.fullPath,
			}

			// 当前节点重新设置属性
			// 原节点的子节点已经被child继承了，新的节点将作为原节点的新子节点
			n.children = []*node{&child}
			// 当前节点路径首字符添加到indices
			// indices将添加所有静态子节点路径的首字符，是为了后面更快查找
			// BytesToString将支持unicode字符
			n.indices = bytesToString([]byte{n.path[i]})
			n.path = path[:i]
			n.handler = nil
			n.wildChild = false
			// fullPath改为前缀部分
			n.fullPath = fullPath[:parentFullPathIndex+i]
		}

		// 最长前缀比路径要短
		// 需要截取路径的后半部分添加到树中
		if i < len(path) {
			// 重新设置path为路径的后半部分
			path = path[i:]
			c := path[0]

			// '/' after param
			if n.nType == param && c == '/' && len(n.children) == 1 {
				parentFullPathIndex += len(n.path)
				n = n.children[0]
				n.priority++
				continue walk
			}

			// 如果indices中能够找到新路径的首字符，将新路径添加到对应的子节点中
			for i, max := 0, len(n.indices); i < max; i++ {
				if c == n.indices[i] {
					parentFullPathIndex += len(n.path)
					// 第i个子节点的路由流量添加了，重新排序子节点，将路由数量多的往前排（路由数量越多，查找路由时命中率越高）
					i = n.incrementChildPrio(i)
					// 设置n为n的第i个子节点，进入新一轮的循环
					n = n.children[i]
					continue walk
				}
			}

			// 如果找不到首字符，将新路径插入当前节点
			if c != ':' && c != '*' && n.nType != catchAll {
				// 首字符添加到indices
				n.indices += bytesToString([]byte{c})
				child := &node{
					fullPath: fullPath,
				}
				n.addChild(child)
				// 对新的子节点切片进行排序
				n.incrementChildPrio(len(n.indices) - 1)
				// 将n设置为新建的子节点
				n = child
			} else if n.wildChild {
				// inserting a wildcard node, need to check if it conflicts with the existing wildcard
				n = n.children[len(n.children)-1]
				n.priority++

				// Check if the wildcard matches
				if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
					// Adding a child to a catchAll is not possible
					n.nType != catchAll &&
					// Check for longer wildcard, e.g. :name and :names
					(len(n.path) >= len(path) || path[len(n.path)] == '/') {
					continue walk
				}

				// Wildcard conflict
				pathSeg := path
				if n.nType != catchAll {
					pathSeg = strings.SplitN(pathSeg, "/", 2)[0]
				}
				prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
				return fmt.Errorf("'%s' in new path '%s' conflicts with existing wildcard '%s' in existing prefix '%s'", pathSeg, fullPath, n.path, prefix)
			}

			// 调用insertChild插入子节点
			n.insertChild(path, fullPath, handler)
			return nil
		}

		// 最长前缀等于路径，当前节点就是路径的节点，添加handler
		// 当前节点已经有非空的handler，说明已经注册过了，重复注册了
		if n.handler != nil {
			return fmt.Errorf("handler are already registered for path '%s'", fullPath)
		}
		n.handler = handler
		n.fullPath = fullPath
		return nil
	}
}

func (n *node) insertChild(path, fullPath string, handler http.HttpHandler) error {
	for {
		// Find prefix until first wildcard
		wildcard, i, valid, sep := findWildcard(path)
		if i < 0 { // No wildcard found
			break
		}
		if !valid {
			return fmt.Errorf("only one wildcard ('*' or ':') per path segment is allowed, has: '%s' in path '%s'", wildcard, fullPath)
		}
		if len(wildcard) < 2 {
			return fmt.Errorf("wildcard (':') must be named with a non-empty name in path '%s'", fullPath)
		}
		if wildcard[0] == ':' && (i <= 0 || path[i-1] != '/') {
			return fmt.Errorf("wildcards ':' must be named in path, but named in domain, '%s'", fullPath)
		}
		if wildcard[0] == '*' && sep == '.' {
			return fmt.Errorf("wildcards '*' must be named in path, but named in domain, '%s'", fullPath)
		}
		if wildcard[0] == ':' {
			if i > 0 {
				// Insert prefix before the current wildcard
				n.path = path[:i]
				n.fullPath = fullPath[:len(fullPath)-len(path)] + n.path
				path = path[i:]
			}

			child := &node{
				nType:    param,
				path:     wildcard,
				fullPath: n.fullPath + wildcard, // fullPath,
			}
			n.addChild(child)
			n.wildChild = true
			n = child
			n.priority++

			// if the path doesn't end with the wildcard, then there
			// will be another subpath starting with '/'
			if len(wildcard) < len(path) {
				path = path[len(wildcard):]
				if path[0] != '*' && path[0] != ':' {
					n.indices += bytesToString([]byte{path[0]})
				}

				child := &node{
					priority: 1,
					fullPath: fullPath,
				}
				n.addChild(child)
				n = child
				continue
			}

			// Otherwise we're done. Insert the handle in the new leaf
			n.handler = handler
			return nil
		}

		if i+len(wildcard) != len(path) {
			return fmt.Errorf("catch-all routes are only allowed at the end of the path in path '%s'", fullPath)
		}

		if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
			pathSeg := strings.SplitN(n.children[0].path, "/", 2)[0]
			return fmt.Errorf("catch-all wildcard '%s' in new path '%s' conflicts with existing path segment '%s' in existing prefix '%s%s'", path, fullPath, pathSeg, n.path, pathSeg)
		}

		// currently fixed width 1 for '/'
		i--
		if path[i] != '/' {
			return fmt.Errorf("no / before catch-all in path '%s'", fullPath)
		}

		n.path = path[:i]

		// First node: catchAll node with empty path
		child := &node{
			wildChild: true,
			nType:     catchAll,
			fullPath:  fullPath,
		}

		n.addChild(child)
		n.indices = string('/')
		n = child
		n.priority++

		// second node: node holding the variable
		child = &node{
			path:     path[i:],
			nType:    catchAll,
			handler:  handler,
			priority: 1,
			fullPath: fullPath,
		}
		n.children = []*node{child}

		return nil

	}
	n.path = path
	n.handler = handler
	n.fullPath = fullPath
	return nil
}

// addChild will add a child node, keeping wildcardChild at the end
func (n *node) addChild(child *node) {
	if n.wildChild && len(n.children) > 0 {
		wildcardChild := n.children[len(n.children)-1]
		n.children = append(n.children[:len(n.children)-1], child, wildcardChild)
	} else {
		n.children = append(n.children, child)
	}
}

func findWildcard(path string) (wildcard string, i int, valid bool, sep byte) {
	// Find start
	for start, c := range []byte(path) {
		// A wildcard starts with ':' (param) or '*' (catch-all)
		if c != ':' && c != '*' {
			continue
		}

		// Find end and check for invalid characters
		valid = true
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '/', domainSep:
				return path[start : start+1+end], start, valid, c
			case ':', '*':
				valid = false
			}
		}
		return path[start:], start, valid, 0
	}
	return "", -1, false, 0
}

func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// Increments priority of the given child and reorders if necessary
func (n *node) incrementChildPrio(pos int) int {
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	// Adjust position (move to front)
	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		// Swap node positions
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}

	// Build new index char string
	if newPos != pos {
		n.indices = n.indices[:newPos] + // Unchanged prefix, might be empty
			n.indices[pos:pos+1] + // The index char we move
			n.indices[newPos:pos] + n.indices[pos+1:] // Rest without char at 'pos'
	}

	return newPos
}

type nodeValue struct {
	handler  http.HttpHandler
	params   *http.Params
	tsr      bool
	pre      bool
	fullPath string // 匹配到的路由全路径
}

type skippedNode struct {
	path        string
	node        *node
	paramsCount int16
}

// type SkippedCatchAllNode struct {
// 	path        string
// 	node        []*node
// 	paramsCount int16
// }

// type SkippedParamsNode struct {
// 	path        string
// 	node        *node
// 	paramsCount int16
// }

func (n *node) getValue(path string, skippedNodes *[]skippedNode, params *http.Params, unescape bool) (value nodeValue, err error) {
	var globalParamsCount int16
walk:
	for {
		prefix := n.path
		if len(path) > len(prefix) {
			if path[:len(prefix)] == prefix {
				path = path[len(prefix):]

				// Try all the non-wildcard children first by matching the indices
				idxc := path[0]
				for i, c := range []byte(n.indices) {
					if c == idxc {
						//  strings.HasPrefix(n.children[len(n.children)-1].path, ":") == n.wildChild
						if n.wildChild {
							if skippedNodes == nil {
								skips := make([]skippedNode, 0, maxSection)
								skippedNodes = &skips
							}

							index := len(*skippedNodes)
							*skippedNodes = (*skippedNodes)[:index+1]
							(*skippedNodes)[index] = skippedNode{
								path: prefix + path,
								node: &node{
									path:      n.path,
									wildChild: n.wildChild,
									nType:     n.nType,
									priority:  n.priority,
									children:  n.children,
									handler:   n.handler,
									fullPath:  n.fullPath,
								},
								paramsCount: globalParamsCount,
							}
						}
						n = n.children[i]
						continue walk
					}
				}

				if !n.wildChild {
					if !value.tsr && !value.pre {
						if value.handler = n.handler; value.handler != nil {
							value.fullPath = n.fullPath
							if path == "/" {
								value.tsr = true
							} else {
								value.pre = true
							}
						}
					}
					// 从paramsSkips遍历，从cathcAll中遍历
					if skippedNodes != nil {
						for length := len(*skippedNodes); length > 0; length-- {
							skippedNode := (*skippedNodes)[length-1]
							*skippedNodes = (*skippedNodes)[:length-1]
							if strings.HasSuffix(skippedNode.path, path) {
								path = skippedNode.path
								n = skippedNode.node
								if value.params != nil {
									*value.params = (*value.params)[:skippedNode.paramsCount]
								}
								globalParamsCount = skippedNode.paramsCount
								continue walk
							}
						}
					}

					return
				}

				// Handle wildcard child, which is always at the end of the array
				n = n.children[len(n.children)-1]
				globalParamsCount++

				switch n.nType {
				case param:
					// fix truncate the parameter
					// tree_test.go  line: 204

					// Find param end (either '/' or path end)
					end := 0
					for end < len(path) && path[end] != '/' {
						end++
					}

					// Save param value
					if params == nil {
						pms := make(http.Params, 0, maxParamsCount)
						params = &pms
					}
					if params != nil && cap(*params) > 0 {
						if value.params == nil {
							value.params = params
						}
						// Expand slice within preallocated capacity
						i := len(*value.params)
						*value.params = (*value.params)[:i+1]
						val := path[:end]
						if unescape {
							if v, err := url.QueryUnescape(val); err == nil {
								val = v
							}
						}
						(*value.params)[i] = http.Param{
							Key:   n.path[1:],
							Value: val,
						}
					}

					// we need to go deeper!
					if end < len(path) {
						if len(n.children) > 0 {
							path = path[end:]
							n = n.children[0]
							continue walk
						}

						// ... but we can't
						value.tsr = len(path) == end+1
						return
					}

					if value.handler = n.handler; value.handler != nil {
						value.fullPath = n.fullPath
						return
					}
					if len(n.children) == 1 {
						// No handle found. Check if a handle for this path + a
						// trailing slash exists for TSR recommendation
						n = n.children[0]
						value.tsr = (n.path == "/" && n.handler != nil) || (n.path == "" && n.indices == "/")
					}
					return

				case catchAll:
					// Save param value
					if params != nil {
						if value.params == nil {
							value.params = params
						}
						// Expand slice within preallocated capacity
						i := len(*value.params)
						*value.params = (*value.params)[:i+1]
						val := path
						if unescape {
							if v, err := url.QueryUnescape(path); err == nil {
								val = v
							}
						}
						(*value.params)[i] = http.Param{
							Key:   n.path[2:],
							Value: val,
						}
					}

					value.handler = n.handler
					value.fullPath = n.fullPath
					return

				default:
					err = fmt.Errorf("invalid node type")
					return
				}
			}
		}

		if path == prefix {
			// If the current path does not equal '/' and the node does not have a registered handle and the most recently matched node has a child node
			// the current node needs to roll back to last valid skippedNode
			if n.handler == nil && path != "/" {
				if skippedNodes != nil {
					for length := len(*skippedNodes); length > 0; length-- {
						skippedNode := (*skippedNodes)[length-1]
						*skippedNodes = (*skippedNodes)[:length-1]
						if strings.HasSuffix(skippedNode.path, path) {
							path = skippedNode.path
							n = skippedNode.node
							if value.params != nil {
								*value.params = (*value.params)[:skippedNode.paramsCount]
							}
							globalParamsCount = skippedNode.paramsCount
							continue walk
						}
					}
				}
				//	n = latestNode.children[len(latestNode.children)-1]
			}
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if value.handler = n.handler; value.handler != nil {
				value.fullPath = n.fullPath
				return
			}

			// If there is no handle for this route, but this route has a
			// wildcard child, there must be a handle for this path with an
			// additional trailing slash
			if path == "/" && n.wildChild && n.nType != root {
				value.tsr = true
				return
			}

			if path == "/" && n.nType == static {
				value.tsr = true
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i, c := range []byte(n.indices) {
				if c == '/' {
					n = n.children[i]
					value.tsr = (len(n.path) == 1 && n.handler != nil) ||
						(n.nType == catchAll && n.children[0].handler != nil)
					return
				}
			}

			return
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		if !value.tsr && (path == "/" ||
			(len(prefix) == len(path)+1 && prefix[len(path)] == '/' && path == prefix[:len(prefix)-1])) {
			if value.handler = n.handler; value.handler != nil {
				value.fullPath = n.fullPath
				value.tsr = true
				value.pre = false
			}
		}

		// roll back to last valid skippedNode
		if !value.tsr && path != "/" {
			if skippedNodes != nil {
				for length := len(*skippedNodes); length > 0; length-- {
					skippedNode := (*skippedNodes)[length-1]
					*skippedNodes = (*skippedNodes)[:length-1]
					if strings.HasSuffix(skippedNode.path, path) {
						path = skippedNode.path
						n = skippedNode.node
						if value.params != nil {
							*value.params = (*value.params)[:skippedNode.paramsCount]
						}
						globalParamsCount = skippedNode.paramsCount
						continue walk
					}
				}
			}
		}
		return
	}
}
