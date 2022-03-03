require 'spec_helper'
require 'open3'

GOLANG_TEST_FAILURE = /FAIL:/

def check_output(output, wait_thr)
  test_failures = []

  output.each_line do |line|
    puts line
    test_failures << line.strip if line =~ GOLANG_TEST_FAILURE
  end

  if test_failures.empty? && !wait_thr.value.success?
    test_failures << "Test command exited with status (#{wait_thr.value.exitstatus}) but no failures were captured."
  end

  test_failures
end

print `cat /etc/os-release`
print `uname -a`

Dir.glob('/tmp/system-probe/head/**/testsuite').each do |f|
  pkg = f.delete_prefix('/tmp/system-probe/head').delete_suffix('/testsuite')
  describe "HEAD system-probe benchmarks for #{pkg}" do
    it 'successfully runs' do
      Dir.chdir(File.dirname(f)) do
        Open3.popen2e({"DD_SYSTEM_PROBE_BPF_DIR"=>"/tmp/system-probe/head/pkg/ebpf/bytecode/build"}, "sudo", "-E", f, "-test.v", "-test.run=^$", "-test.benchmem", "-test.bench=.") do |_, output, wait_thr|
          test_failures = check_output(output, wait_thr)
          expect(test_failures).to be_empty, test_failures.join("\n")
        end
      end
    end
  end
end

Dir.glob('/tmp/system-probe/main/**/testsuite').each do |f|
  pkg = f.delete_prefix('/tmp/system-probe/main').delete_suffix('/testsuite')
  describe "MAIN system-probe benchmarks for #{pkg}" do
    it 'successfully runs' do
      Dir.chdir(File.dirname(f)) do
        Open3.popen2e({"DD_SYSTEM_PROBE_BPF_DIR"=>"/tmp/system-probe/main/pkg/ebpf/bytecode/build"}, "sudo", "-E", f, "-test.v", "-test.run=^$", "-test.benchmem", "-test.bench .") do |_, output, wait_thr|
          test_failures = check_output(output, wait_thr)
          expect(test_failures).to be_empty, test_failures.join("\n")
        end
      end
    end
  end
end
