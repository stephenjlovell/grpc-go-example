# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: greet.proto

require 'google/protobuf'

Google::Protobuf::DescriptorPool.generated_pool.build do
  add_file("greet.proto", :syntax => :proto3) do
    add_message "greet.Greeting" do
      optional :first_name, :string, 1
      optional :last_name, :string, 2
    end
    add_message "greet.GreetRequest" do
      optional :greeting, :message, 1, "greet.Greeting"
    end
    add_message "greet.GreetResponse" do
      optional :response, :string, 1
    end
    add_message "greet.GreetManyTimesRequest" do
      optional :greeting, :message, 1, "greet.Greeting"
    end
    add_message "greet.GreetManyTimesResponse" do
      optional :response, :string, 1
    end
    add_message "greet.LongGreetRequest" do
      optional :greeting, :message, 1, "greet.Greeting"
    end
    add_message "greet.LongGreetResponse" do
      optional :response, :string, 1
    end
  end
end

module Greet
  Greeting = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("greet.Greeting").msgclass
  GreetRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("greet.GreetRequest").msgclass
  GreetResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("greet.GreetResponse").msgclass
  GreetManyTimesRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("greet.GreetManyTimesRequest").msgclass
  GreetManyTimesResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("greet.GreetManyTimesResponse").msgclass
  LongGreetRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("greet.LongGreetRequest").msgclass
  LongGreetResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("greet.LongGreetResponse").msgclass
end
