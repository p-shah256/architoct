// Create application user
db.createUser({
  user: 'mvp_user',
  pwd: 'mvp_pwd',  // Change this in production
  roles: [{
    role: 'readWrite',
    db: 'mvp_db'
  }]
});

// Switch to forum database
db = db.getSiblingDB('mvp_db');

// Create collections with schema validation
db.createCollection('users', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['_id', 'created_at'], // _id = fingreprint.js atm
      properties: {
        _id: { bsonType: 'string' },
        created_at: { bsonType: 'timestamp' },
        last_login: { bsonType: 'timestamp' }
      }
    }
  }
});

db.createCollection('stories', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      // NOTE: post_id is generated by backend
      required: ['user_id', 'title', 'body', '_id', 'created_at'],
      properties: {
        // STORY COMPONENTS --------------------------------
        _id: { bsonType: 'string' }, //nanoID maybe with 4 chars
        user_id: { bsonType: 'string' },
        title: { bsonType: 'string' },
        body: { bsonType: 'string' },
        created_at: { bsonType: 'date' },
        // METADATA (needs update frequently) -------------
        upvote_count: { bsonType: 'int' },
        upvoted_by: { // user_ids
          bsonType: 'array', // hashbased structure so retreival is O(1)
          items: { bsonType: 'string' }  // user_ids
        },
        reply_count: { bsonType: 'int' },
        replies: {
          bsonType: 'array',
          items: { bsonType: 'objectId' } // comment_ids
        }
      }
    }
  }
});

db.createCollection('comments', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['user_id', 'body', 'post_id', 'created_at'], //_id will be auto
      properties: {
        // COMMENT DATA -----------------------------------
        // _id: {bsonType: 'string'}, // lets use auto
        post_id: { bsonType: 'string' }, //nanoID maybe with 4 chars
        user_id: { bsonType: 'string' },
        body: { bsonType: 'string' },
        created_at: { bsonType: 'date' },
        is_deleted: { bsonType: 'bool' },
        // META DATA(needs update frequently) --------------
        upvote_count: { bsonType: 'int' },
        upvoted_by: { // user_ids
          bsonType: 'array', // hashbased structure so retreival is O(1)
          items: { bsonType: 'string' }  // user_ids
        },
        reply_count: { bsonType: 'int' },
        replies: { // stores comment ids here..
          bsonType: 'array',
          items: {bsonType: 'string'}
        }
      }
    }
  }
});

// Create indexes
// 1. get top x posts in x days
db.stories.createIndex({ created_at: -1, upvote_count: -1 });
// 2. get top x comments for x posts
db.comments.createIndex({ post_id: 1, upvote_count: -1 });
