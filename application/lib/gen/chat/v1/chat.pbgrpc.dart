// This is a generated file - do not edit.
//
// Generated from chat/v1/chat.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'chat.pb.dart' as $0;

export 'chat.pb.dart';

@$pb.GrpcServiceName('chat.v1.ChatService')
class ChatServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  ChatServiceClient(super.channel, {super.options, super.interceptors});

  /// Registers a new participant with a client-chosen nickname and color.
  $grpc.ResponseFuture<$0.JoinResponse> join(
    $0.JoinRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$join, request, options: options);
  }

  /// Sends a plain-text message to the room.
  $grpc.ResponseFuture<$0.SendMessageResponse> sendMessage(
    $0.SendMessageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendMessage, request, options: options);
  }

  /// Streams every message posted to the room, from the moment of
  /// subscribing onward.
  $grpc.ResponseStream<$0.ChatMessage> subscribe(
    $0.SubscribeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$subscribe, $async.Stream.fromIterable([request]),
        options: options);
  }

  // method descriptors

  static final _$join = $grpc.ClientMethod<$0.JoinRequest, $0.JoinResponse>(
      '/chat.v1.ChatService/Join',
      ($0.JoinRequest value) => value.writeToBuffer(),
      $0.JoinResponse.fromBuffer);
  static final _$sendMessage =
      $grpc.ClientMethod<$0.SendMessageRequest, $0.SendMessageResponse>(
          '/chat.v1.ChatService/SendMessage',
          ($0.SendMessageRequest value) => value.writeToBuffer(),
          $0.SendMessageResponse.fromBuffer);
  static final _$subscribe =
      $grpc.ClientMethod<$0.SubscribeRequest, $0.ChatMessage>(
          '/chat.v1.ChatService/Subscribe',
          ($0.SubscribeRequest value) => value.writeToBuffer(),
          $0.ChatMessage.fromBuffer);
}

@$pb.GrpcServiceName('chat.v1.ChatService')
abstract class ChatServiceBase extends $grpc.Service {
  $core.String get $name => 'chat.v1.ChatService';

  ChatServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.JoinRequest, $0.JoinResponse>(
        'Join',
        join_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.JoinRequest.fromBuffer(value),
        ($0.JoinResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.SendMessageRequest, $0.SendMessageResponse>(
            'SendMessage',
            sendMessage_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.SendMessageRequest.fromBuffer(value),
            ($0.SendMessageResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SubscribeRequest, $0.ChatMessage>(
        'Subscribe',
        subscribe_Pre,
        false,
        true,
        ($core.List<$core.int> value) => $0.SubscribeRequest.fromBuffer(value),
        ($0.ChatMessage value) => value.writeToBuffer()));
  }

  $async.Future<$0.JoinResponse> join_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.JoinRequest> $request) async {
    return join($call, await $request);
  }

  $async.Future<$0.JoinResponse> join(
      $grpc.ServiceCall call, $0.JoinRequest request);

  $async.Future<$0.SendMessageResponse> sendMessage_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SendMessageRequest> $request) async {
    return sendMessage($call, await $request);
  }

  $async.Future<$0.SendMessageResponse> sendMessage(
      $grpc.ServiceCall call, $0.SendMessageRequest request);

  $async.Stream<$0.ChatMessage> subscribe_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SubscribeRequest> $request) async* {
    yield* subscribe($call, await $request);
  }

  $async.Stream<$0.ChatMessage> subscribe(
      $grpc.ServiceCall call, $0.SubscribeRequest request);
}
