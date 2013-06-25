## Supervisord

[Supervisord](http://supervisord.org/) will make sure your web app is always up.

1. Installation

		wget http://pypi.python.org/packages/2.7/s/setuptools/setuptools-0.6c11-py2.7.egg
		
		sh setuptools-0.6c11-py2.7.egg
		
		easy_install supervisor
		
		echo_supervisord_conf >/etc/supervisord.conf
		
		mkdir /etc/supervisord.conf.d

2. Configure `/etc/supervisord.conf`

		[include]
		files = /etc/supervisord.conf.d/*.conf

3. Add new application

		cd /etc/supervisord.conf.d
		vim beepkg.conf
	
	Configuration file:
	
		[program:beepkg]
		directory = /opt/app/beepkg
		command = /opt/app/beepkg/beepkg
		autostart = true
		startsecs = 5
		user = root
		redirect_stderr = true
		stdout_logfile = /var/log/supervisord/beepkg.log
