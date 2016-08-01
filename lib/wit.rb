require 'uri'
require 'json'

module Oncall
  module Wit
    #included(self) do
    #  puts "Create wit client"
    #end

    def self.wit
      @wit ||= Client.new ENV["WIT_ACCESS_TOKEN"]
    end

    class Client
      include HTTParty
      base_uri 'https://api.wit.ai'

      def initialize(token)
        @token = token
      end

      def parse(message)
        response = request(message)
        intent   = process_response(response)
        intent
      end

      private
      def request(q)
        q = { query: {
              v: '20160731',
              q: URI.escape(q)
            }}

        self.class.get("/message", q)
        response.body
      end

      def process_response(response)
        Response.new response
      end
    end

    class Response
      def initialize(body)
        @body = body
        @response = JSON.parse(body)
      end

      def type
        "attribute"
      end

      def [](y)
        @response["entities"][y]
      end
    end

  end
end
