namespace :deploy do
  desc "build"
  task :build do
    sh("go download")
    sh("build/build.sh")
  end

  desc "restart service"
  task :restart do
    invoke "deploy:stop"
    sleep(3)
    invoke "deploy:start"
  end

  desc "systemd stop"
  task :"sys-stop" do
    app = application
    on roles(:app) do
      puts("stopping...")
      execute("sudo systemctl stop #{app}")
    end
  end

  desc "systemd start"
  task :"sys-start" do
    app = application
    on roles(:app) do
      puts("starting...")
      execute("sudo systemctl start #{app}")
    end
  end

  desc "systemd restart"
  task :"sys-restart" do
    app = application
    on roles(:app) do
      puts("restarting...")
      invoke "deploy:sys-stop"
      sleep(3)
      invoke "deploy:sys-start"
    end
  end
end
