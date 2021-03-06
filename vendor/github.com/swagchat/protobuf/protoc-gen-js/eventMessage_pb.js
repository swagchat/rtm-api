/**
 * @fileoverview
 * @enhanceable
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

var gogoproto_gogo_pb = require('./gogoproto/gogo_pb.js');
var roomMessage_pb = require('./roomMessage_pb.js');
goog.exportSymbol('proto.swagchat.protobuf.EventData', null, global);
goog.exportSymbol('proto.swagchat.protobuf.EventType', null, global);
goog.exportSymbol('proto.swagchat.protobuf.RoomEventPayload', null, global);

/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.swagchat.protobuf.EventData = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.swagchat.protobuf.EventData.repeatedFields_, null);
};
goog.inherits(proto.swagchat.protobuf.EventData, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.swagchat.protobuf.EventData.displayName = 'proto.swagchat.protobuf.EventData';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.swagchat.protobuf.EventData.repeatedFields_ = [3];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.swagchat.protobuf.EventData.prototype.toObject = function(opt_includeInstance) {
  return proto.swagchat.protobuf.EventData.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.swagchat.protobuf.EventData} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.swagchat.protobuf.EventData.toObject = function(includeInstance, msg) {
  var f, obj = {
    type: jspb.Message.getField(msg, 1),
    data: msg.getData_asB64(),
    userIdsList: jspb.Message.getRepeatedField(msg, 3)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.swagchat.protobuf.EventData}
 */
proto.swagchat.protobuf.EventData.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.swagchat.protobuf.EventData;
  return proto.swagchat.protobuf.EventData.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.swagchat.protobuf.EventData} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.swagchat.protobuf.EventData}
 */
proto.swagchat.protobuf.EventData.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {!proto.swagchat.protobuf.EventType} */ (reader.readEnum());
      msg.setType(value);
      break;
    case 2:
      var value = /** @type {!Uint8Array} */ (reader.readBytes());
      msg.setData(value);
      break;
    case 3:
      var value = /** @type {string} */ (reader.readString());
      msg.addUserIds(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.swagchat.protobuf.EventData.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.swagchat.protobuf.EventData.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.swagchat.protobuf.EventData} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.swagchat.protobuf.EventData.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {!proto.swagchat.protobuf.EventType} */ (jspb.Message.getField(message, 1));
  if (f != null) {
    writer.writeEnum(
      1,
      f
    );
  }
  f = /** @type {!(string|Uint8Array)} */ (jspb.Message.getField(message, 2));
  if (f != null) {
    writer.writeBytes(
      2,
      f
    );
  }
  f = message.getUserIdsList();
  if (f.length > 0) {
    writer.writeRepeatedString(
      3,
      f
    );
  }
};


/**
 * optional EventType type = 1;
 * @return {!proto.swagchat.protobuf.EventType}
 */
proto.swagchat.protobuf.EventData.prototype.getType = function() {
  return /** @type {!proto.swagchat.protobuf.EventType} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/** @param {!proto.swagchat.protobuf.EventType} value */
proto.swagchat.protobuf.EventData.prototype.setType = function(value) {
  jspb.Message.setField(this, 1, value);
};


proto.swagchat.protobuf.EventData.prototype.clearType = function() {
  jspb.Message.setField(this, 1, undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.swagchat.protobuf.EventData.prototype.hasType = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * optional bytes data = 2;
 * @return {!(string|Uint8Array)}
 */
proto.swagchat.protobuf.EventData.prototype.getData = function() {
  return /** @type {!(string|Uint8Array)} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * optional bytes data = 2;
 * This is a type-conversion wrapper around `getData()`
 * @return {string}
 */
proto.swagchat.protobuf.EventData.prototype.getData_asB64 = function() {
  return /** @type {string} */ (jspb.Message.bytesAsB64(
      this.getData()));
};


/**
 * optional bytes data = 2;
 * Note that Uint8Array is not supported on all browsers.
 * @see http://caniuse.com/Uint8Array
 * This is a type-conversion wrapper around `getData()`
 * @return {!Uint8Array}
 */
proto.swagchat.protobuf.EventData.prototype.getData_asU8 = function() {
  return /** @type {!Uint8Array} */ (jspb.Message.bytesAsU8(
      this.getData()));
};


/** @param {!(string|Uint8Array)} value */
proto.swagchat.protobuf.EventData.prototype.setData = function(value) {
  jspb.Message.setField(this, 2, value);
};


proto.swagchat.protobuf.EventData.prototype.clearData = function() {
  jspb.Message.setField(this, 2, undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.swagchat.protobuf.EventData.prototype.hasData = function() {
  return jspb.Message.getField(this, 2) != null;
};


/**
 * repeated string user_ids = 3;
 * @return {!Array.<string>}
 */
proto.swagchat.protobuf.EventData.prototype.getUserIdsList = function() {
  return /** @type {!Array.<string>} */ (jspb.Message.getRepeatedField(this, 3));
};


/** @param {!Array.<string>} value */
proto.swagchat.protobuf.EventData.prototype.setUserIdsList = function(value) {
  jspb.Message.setField(this, 3, value || []);
};


/**
 * @param {!string} value
 * @param {number=} opt_index
 */
proto.swagchat.protobuf.EventData.prototype.addUserIds = function(value, opt_index) {
  jspb.Message.addToRepeatedField(this, 3, value, opt_index);
};


proto.swagchat.protobuf.EventData.prototype.clearUserIdsList = function() {
  this.setUserIdsList([]);
};



/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.swagchat.protobuf.RoomEventPayload = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.swagchat.protobuf.RoomEventPayload.repeatedFields_, null);
};
goog.inherits(proto.swagchat.protobuf.RoomEventPayload, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  proto.swagchat.protobuf.RoomEventPayload.displayName = 'proto.swagchat.protobuf.RoomEventPayload';
}
/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.swagchat.protobuf.RoomEventPayload.repeatedFields_ = [2];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto suitable for use in Soy templates.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     com.google.apps.jspb.JsClassTemplate.JS_RESERVED_WORDS.
 * @param {boolean=} opt_includeInstance Whether to include the JSPB instance
 *     for transitional soy proto support: http://goto/soy-param-migration
 * @return {!Object}
 */
proto.swagchat.protobuf.RoomEventPayload.prototype.toObject = function(opt_includeInstance) {
  return proto.swagchat.protobuf.RoomEventPayload.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Whether to include the JSPB
 *     instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.swagchat.protobuf.RoomEventPayload} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.swagchat.protobuf.RoomEventPayload.toObject = function(includeInstance, msg) {
  var f, obj = {
    roomId: jspb.Message.getField(msg, 1),
    usersList: jspb.Message.toObjectList(msg.getUsersList(),
    roomMessage_pb.MiniUser.toObject, includeInstance)
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.swagchat.protobuf.RoomEventPayload}
 */
proto.swagchat.protobuf.RoomEventPayload.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.swagchat.protobuf.RoomEventPayload;
  return proto.swagchat.protobuf.RoomEventPayload.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.swagchat.protobuf.RoomEventPayload} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.swagchat.protobuf.RoomEventPayload}
 */
proto.swagchat.protobuf.RoomEventPayload.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {string} */ (reader.readString());
      msg.setRoomId(value);
      break;
    case 2:
      var value = new roomMessage_pb.MiniUser;
      reader.readMessage(value,roomMessage_pb.MiniUser.deserializeBinaryFromReader);
      msg.addUsers(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.swagchat.protobuf.RoomEventPayload.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.swagchat.protobuf.RoomEventPayload.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.swagchat.protobuf.RoomEventPayload} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.swagchat.protobuf.RoomEventPayload.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = /** @type {string} */ (jspb.Message.getField(message, 1));
  if (f != null) {
    writer.writeString(
      1,
      f
    );
  }
  f = message.getUsersList();
  if (f.length > 0) {
    writer.writeRepeatedMessage(
      2,
      f,
      roomMessage_pb.MiniUser.serializeBinaryToWriter
    );
  }
};


/**
 * optional string room_id = 1;
 * @return {string}
 */
proto.swagchat.protobuf.RoomEventPayload.prototype.getRoomId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/** @param {string} value */
proto.swagchat.protobuf.RoomEventPayload.prototype.setRoomId = function(value) {
  jspb.Message.setField(this, 1, value);
};


proto.swagchat.protobuf.RoomEventPayload.prototype.clearRoomId = function() {
  jspb.Message.setField(this, 1, undefined);
};


/**
 * Returns whether this field is set.
 * @return {!boolean}
 */
proto.swagchat.protobuf.RoomEventPayload.prototype.hasRoomId = function() {
  return jspb.Message.getField(this, 1) != null;
};


/**
 * repeated MiniUser users = 2;
 * @return {!Array.<!proto.swagchat.protobuf.MiniUser>}
 */
proto.swagchat.protobuf.RoomEventPayload.prototype.getUsersList = function() {
  return /** @type{!Array.<!proto.swagchat.protobuf.MiniUser>} */ (
    jspb.Message.getRepeatedWrapperField(this, roomMessage_pb.MiniUser, 2));
};


/** @param {!Array.<!proto.swagchat.protobuf.MiniUser>} value */
proto.swagchat.protobuf.RoomEventPayload.prototype.setUsersList = function(value) {
  jspb.Message.setRepeatedWrapperField(this, 2, value);
};


/**
 * @param {!proto.swagchat.protobuf.MiniUser=} opt_value
 * @param {number=} opt_index
 * @return {!proto.swagchat.protobuf.MiniUser}
 */
proto.swagchat.protobuf.RoomEventPayload.prototype.addUsers = function(opt_value, opt_index) {
  return jspb.Message.addToRepeatedWrapperField(this, 2, opt_value, proto.swagchat.protobuf.MiniUser, opt_index);
};


proto.swagchat.protobuf.RoomEventPayload.prototype.clearUsersList = function() {
  this.setUsersList([]);
};


/**
 * @enum {number}
 */
proto.swagchat.protobuf.EventType = {
  EMPTYEVENT: 0,
  MESSAGEEVENT: 1,
  ROOMEVENT: 2
};

goog.object.extend(exports, proto.swagchat.protobuf);
