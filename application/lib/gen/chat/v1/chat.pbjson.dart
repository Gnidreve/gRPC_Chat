// This is a generated file - do not edit.
//
// Generated from chat/v1/chat.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use userDescriptor instead')
const User$json = {
  '1': 'User',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'nickname', '3': 2, '4': 1, '5': 9, '10': 'nickname'},
    {'1': 'color', '3': 3, '4': 1, '5': 9, '10': 'color'},
  ],
};

/// Descriptor for `User`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userDescriptor = $convert.base64Decode(
    'CgRVc2VyEg4KAmlkGAEgASgJUgJpZBIaCghuaWNrbmFtZRgCIAEoCVIIbmlja25hbWUSFAoFY2'
    '9sb3IYAyABKAlSBWNvbG9y');

@$core.Deprecated('Use chatMessageDescriptor instead')
const ChatMessage$json = {
  '1': 'ChatMessage',
  '2': [
    {'1': 'user', '3': 1, '4': 1, '5': 11, '6': '.chat.v1.User', '10': 'user'},
    {'1': 'text', '3': 2, '4': 1, '5': 9, '10': 'text'},
    {
      '1': 'sent_at',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'sentAt'
    },
    {
      '1': 'distance_km',
      '3': 4,
      '4': 1,
      '5': 5,
      '9': 0,
      '10': 'distanceKm',
      '17': true
    },
  ],
  '8': [
    {'1': '_distance_km'},
  ],
};

/// Descriptor for `ChatMessage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatMessageDescriptor = $convert.base64Decode(
    'CgtDaGF0TWVzc2FnZRIhCgR1c2VyGAEgASgLMg0uY2hhdC52MS5Vc2VyUgR1c2VyEhIKBHRleH'
    'QYAiABKAlSBHRleHQSMwoHc2VudF9hdBgDIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3Rh'
    'bXBSBnNlbnRBdBIkCgtkaXN0YW5jZV9rbRgEIAEoBUgAUgpkaXN0YW5jZUttiAEBQg4KDF9kaX'
    'N0YW5jZV9rbQ==');

@$core.Deprecated('Use joinRequestDescriptor instead')
const JoinRequest$json = {
  '1': 'JoinRequest',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'nickname', '3': 2, '4': 1, '5': 9, '10': 'nickname'},
    {'1': 'color', '3': 3, '4': 1, '5': 9, '10': 'color'},
  ],
};

/// Descriptor for `JoinRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List joinRequestDescriptor = $convert.base64Decode(
    'CgtKb2luUmVxdWVzdBIOCgJpZBgBIAEoCVICaWQSGgoIbmlja25hbWUYAiABKAlSCG5pY2tuYW'
    '1lEhQKBWNvbG9yGAMgASgJUgVjb2xvcg==');

@$core.Deprecated('Use joinResponseDescriptor instead')
const JoinResponse$json = {
  '1': 'JoinResponse',
  '2': [
    {'1': 'user', '3': 1, '4': 1, '5': 11, '6': '.chat.v1.User', '10': 'user'},
  ],
};

/// Descriptor for `JoinResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List joinResponseDescriptor = $convert.base64Decode(
    'CgxKb2luUmVzcG9uc2USIQoEdXNlchgBIAEoCzINLmNoYXQudjEuVXNlclIEdXNlcg==');

@$core.Deprecated('Use sendMessageRequestDescriptor instead')
const SendMessageRequest$json = {
  '1': 'SendMessageRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'text', '3': 2, '4': 1, '5': 9, '10': 'text'},
    {'1': 'lat', '3': 3, '4': 1, '5': 1, '9': 0, '10': 'lat', '17': true},
    {'1': 'lng', '3': 4, '4': 1, '5': 1, '9': 1, '10': 'lng', '17': true},
  ],
  '8': [
    {'1': '_lat'},
    {'1': '_lng'},
  ],
};

/// Descriptor for `SendMessageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendMessageRequestDescriptor = $convert.base64Decode(
    'ChJTZW5kTWVzc2FnZVJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEhIKBHRleHQYAi'
    'ABKAlSBHRleHQSFQoDbGF0GAMgASgBSABSA2xhdIgBARIVCgNsbmcYBCABKAFIAVIDbG5niAEB'
    'QgYKBF9sYXRCBgoEX2xuZw==');

@$core.Deprecated('Use sendMessageResponseDescriptor instead')
const SendMessageResponse$json = {
  '1': 'SendMessageResponse',
  '2': [
    {
      '1': 'message',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.chat.v1.ChatMessage',
      '10': 'message'
    },
  ],
};

/// Descriptor for `SendMessageResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendMessageResponseDescriptor = $convert.base64Decode(
    'ChNTZW5kTWVzc2FnZVJlc3BvbnNlEi4KB21lc3NhZ2UYASABKAsyFC5jaGF0LnYxLkNoYXRNZX'
    'NzYWdlUgdtZXNzYWdl');

@$core.Deprecated('Use subscribeRequestDescriptor instead')
const SubscribeRequest$json = {
  '1': 'SubscribeRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'lat', '3': 2, '4': 1, '5': 1, '9': 0, '10': 'lat', '17': true},
    {'1': 'lng', '3': 3, '4': 1, '5': 1, '9': 1, '10': 'lng', '17': true},
  ],
  '8': [
    {'1': '_lat'},
    {'1': '_lng'},
  ],
};

/// Descriptor for `SubscribeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscribeRequestDescriptor = $convert.base64Decode(
    'ChBTdWJzY3JpYmVSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIVCgNsYXQYAiABKA'
    'FIAFIDbGF0iAEBEhUKA2xuZxgDIAEoAUgBUgNsbmeIAQFCBgoEX2xhdEIGCgRfbG5n');

@$core.Deprecated('Use presenceUpdateDescriptor instead')
const PresenceUpdate$json = {
  '1': 'PresenceUpdate',
  '2': [
    {'1': 'online_count', '3': 1, '4': 1, '5': 5, '10': 'onlineCount'},
  ],
};

/// Descriptor for `PresenceUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List presenceUpdateDescriptor = $convert.base64Decode(
    'Cg5QcmVzZW5jZVVwZGF0ZRIhCgxvbmxpbmVfY291bnQYASABKAVSC29ubGluZUNvdW50');

@$core.Deprecated('Use chatEventDescriptor instead')
const ChatEvent$json = {
  '1': 'ChatEvent',
  '2': [
    {
      '1': 'message',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.chat.v1.ChatMessage',
      '9': 0,
      '10': 'message'
    },
    {
      '1': 'presence',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.chat.v1.PresenceUpdate',
      '9': 0,
      '10': 'presence'
    },
  ],
  '8': [
    {'1': 'event'},
  ],
};

/// Descriptor for `ChatEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List chatEventDescriptor = $convert.base64Decode(
    'CglDaGF0RXZlbnQSMAoHbWVzc2FnZRgBIAEoCzIULmNoYXQudjEuQ2hhdE1lc3NhZ2VIAFIHbW'
    'Vzc2FnZRI1CghwcmVzZW5jZRgCIAEoCzIXLmNoYXQudjEuUHJlc2VuY2VVcGRhdGVIAFIIcHJl'
    'c2VuY2VCBwoFZXZlbnQ=');
