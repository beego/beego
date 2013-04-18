## supervisord安装

1. setuptools安装

		wget http://pypi.python.org/packages/2.7/s/setuptools/setuptools-0.6c11-py2.7.egg
 
		sh setuptools-0.6c11-py2.7.egg 
 
		easy_install supervisor
 		
		echo_supervisord_conf >/etc/supervisord.conf
 
		mkdir /etc/supervisord.conf.d
 
2. 修改配置/etc/supervisord.conf 

		[include]
		files = /etc/supervisord.conf.d/*.conf
 
3. 新建管理的应用

		cd /etc/supervisord.conf.d
		vim ddq.conf
	
	配置文件：
	
		[program:ddq]
		directory = /opt/app/ddq
		command = /opt/app/ddq/ddq
		autostart = true
		startsecs = 5
		user = root
		redirect_stderr = true
		stdout_logfile = /var/log/supervisord/shorturl.log 