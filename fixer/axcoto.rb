module Oncall
  module Fixer

    class Axcoto
      include Oncall::Router
      include Oncall::Executor

      # Define how to SSH Into server
      ssh do
        {host: "axcoto.com", username: "kurei"}
      end

      route do
        command 'mem' do |data, match|
          say "I'm processing #{match[:some_argument]}"

          result = run("this command is run over SSH connection above")
          output.append result
          # Or we can do some kind of processing
          post_result = crazy_processing result
          # Then we can continue to run another command 
          result = run("another command based on post_result")
          output.append result
        end

        match /axcoto mem/ do |data, match|
          say 'I got'
        end
      end


    end

  end
end
