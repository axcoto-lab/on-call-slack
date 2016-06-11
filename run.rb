require 'slack-ruby-bot'

require_relative 'oncall/oncall'

require_relative 'fixer/axcoto'

module Oncall
  VERSION = '0.1'

  # Fixer list
  include Fixer::Axcoto
  #include Fixer::OtherProvider

end

Oncall::Bot.run
