$: << File.expand_path('../../deployments', __FILE__)
# config valid only for current version of Capistrano
# lock '3.7.1'

set :application, 'devops'

# Default branch is :master
# ask :branch, `git rev-parse --abbrev-ref HEAD`.chomp

# Default deploy_to directory is /var/www/my_app_name
# set :deploy_to, '/var/www/api'
set :deploy_to, '/apps/www'

# Default value for :scm is :git
set :scm, :action

# Default value for :format is :pretty
# set :format, :pretty

# Default value for :log_level is :debug
# set :log_level, :info

# Default value for :pty is false
# set :pty, true

# Default value for :linked_files is []
# set :linked_files, fetch(:linked_files, []).push('config/production.conf')

# Default value for linked_dirs is []
set :linked_dirs, fetch(:linked_dirs, []).push('logs', 'configs')
# Default value for default_env is {}
# set :default_env, { path: "/home/apps/.nvm/v0.12.7/bin:$PATH" }

# Default value for keep_releases is 5
set :keep_releases, 5

def self.application
  fetch(:application)
end