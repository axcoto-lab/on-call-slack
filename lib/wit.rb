require 'uri'
require 'json'

module Oncall
  module Wit

    #def self.included(base)
    #  puts "Create wit client"
    #end

    def wit
      @wit ||= Client.new ENV["WIT_ACCESS_TOKEN"]
    end

    class Client
      include HTTParty
      base_uri 'https://api.wit.ai'

      def initialize(token)
        @token = token
        self.class.headers 'Authorization' => "Bearer #{@token}"
      end

      def parse(message)
        response = request(message)
        intent   = process_response(response)
        intent
      end

      private
      def request(q)
        q = { 
          query: {
              v: '20160731',
              q: q
          },
        }

        response = self.class.get("/message", q)
        response.body
      end

      def process_response(response)
        Response.new response
      end
    end

    class Response
      def initialize(body)
        @body = body
        puts @body
        @response = JSON.parse(body)
      end

      def type
        "attribute"
      end

      def [](y)
        @response["entities"][y].first["value"] unless @response["entities"][y].nil?
      end

      def to_a
        @response["entities"]
      end

      def to_s
        @response.to_s
      end
    end

  end
end
