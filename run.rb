require 'slack-ruby-bot'

require_relative 'oncall/oncall'
require_relative 'fixer/axcoto/base'

module Oncall
  VERSION = '0.1'

  # Fixer list
  include Fixer::Axcoto
  #include Fixer::OtherProvider

  def self.run!
    Oncall::Bot.run
  end
end

Oncall.run!
