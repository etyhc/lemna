// Package arpc 是代理服务器rpc定义
//              定义了2个grpc服务CAgent和SAgent，分别用于客户端和服务器通信
//              客户端<---------->代理<---------->服务器
//                       CAgent          SAgent
//              包含了：
//              客户端rpc(go)封装，开发客户端时使用，亦可自行实现
//              服务器rpc(go)封装，开发服务器时使用，其他语言请自行实现
package arpc
