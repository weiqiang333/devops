namespace :action do

  desc "sync"
  task :sync do
    on roles(:app) do
      execute :mkdir, "-p", release_path
      execute :mkdir, "-p", "#{release_path}/web"
      execute :mkdir, "-p", "#{release_path}/bin"
      upload!("web/static", "#{release_path}/web/", recursive: true)
      upload!("web/templates", "#{release_path}/web/", recursive: true)
      upload!("devops", release_path.join("bin/"))
      upload!("devops-cron", release_path.join("bin/"))
    end
  end

  task :create_release => :sync
  task :check
  task :set_current_revision
end
