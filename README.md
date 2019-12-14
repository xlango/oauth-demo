"# oauth-demo" 
<h3>用户认证服务：</h3>
1.创建用户

	（1）创建用户名密码
	（2）kong创建consumer并创建oauth生成所需client_id、client_secret，添加consumer对应服务插件，例如：限流
	（3）将username和client_id、client_secret的对应关系保存
	
2.web用户登录（用户名密码登录）

	（1）用户名密码验证通过，则通过username查出对应client_id、client_secret
	（2）使用client_id请求kong接口 https://192.168.10.33:8443/server/oauth2/authorize获取 code
	（3）使用code请求kong接口https://192.168.10.33:8443/server/oauth2/token获取access_token和refresh_token
	（4）将access_token和username绑定保存至redis并设置过期时间
	
3.第三方接口

	（1）第三方使用client_id、client_secret获取到access_token和refresh_token
	（2）通过client_id获取到对应username并将access_token和username绑定保存至redis并设置过期时间

<h3>权限服务：</h3>
1.提供远程调用，其他模块可通过传递token获取到username和该用户对应的权限     


<h3>web api服务</h3>
1.api请求Filter中获取出token，将token发送至权限服务获取出用户信息和访问该接口权限     


<h3>api gateway：kong</h3>

    1.配置路由
    2.创建consumer对应用户
    3.oauth2认证，token均由kong生成
    4.针对不同consumer进行限流，就相当于对不同用户进行限流
    5.熔断
    
<h3>服务架构图</h3>
