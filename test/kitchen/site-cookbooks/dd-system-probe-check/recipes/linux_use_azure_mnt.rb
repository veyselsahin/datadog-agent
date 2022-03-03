return unless azure? && !platform?('windows')

mnt_path = ::File.exist?('/mnt/resource') ? '/mnt/resource' : '/mnt'

script 'use mnt' do
  interpreter "bash"
  code <<-EOH
    mkdir -p #{mnt_path}/system-probe
    chmod 0777 #{mnt_path}/system-probe
    ln -s #{mnt_path}/system-probe /tmp/system-probe
  EOH
end
