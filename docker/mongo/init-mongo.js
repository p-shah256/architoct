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
      required: ['_id', 'created_at'], // _id = fingreprint uuid atm
      properties: {                    // and then when they create a username, replace it
        _id: { bsonType: 'string' },
        created_at: { bsonType: 'date' },
        last_login: { bsonType: 'date' }
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
        upvote_count: {
          bsonType: 'int',
        },
        upvoted_by: {
          bsonType: 'array',
          uniqueItems: true,  // Prevent duplicate votes
          items: { bsonType: 'string' }
        },
        reply_count: {
          bsonType: 'int',
          minimum: 0,
        },
        // this just creates a breaking point as now you need multidoc transaction
        // AND harder to keep it consistent
        // AND makes sorting by upvotes more harder
        //
        // replies: {
        //   bsonType: 'array',
        //   items: { bsonType: 'objectId' } // comment_ids
        // }
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
        _id: {bsonType: 'objectId'},
        post_id: { bsonType: 'string' }, //nanoID maybe with 4 chars
        user_id: { bsonType: 'string' },
        body: { bsonType: 'string' },
        created_at: { bsonType: 'date' },
        is_deleted: { bsonType: 'bool' },
        // META DATA(needs update frequently) --------------
        upvote_count: {
          bsonType: 'int',
          minimum: 0,
        },
        upvoted_by: {
          bsonType: 'array',
          uniqueItems: true,
          items: { bsonType: 'string' }
        },
        reply_count: { bsonType: 'int' },
        replies: { // stores comment ids here..
          bsonType: 'array',
          items: {
            bsonType: 'objectId',
          }
        }
      }
    }
  }
});

//STORY INDICES////////////////////////////////////////////////////////////////
// 1. get top posts sorted by upvotes from 'y' days
db.stories.createIndex({
  created_at: -1,        // Primary sort by creation time
  upvote_count: -1,      // Secondary sort by votes
  _id: 1                 // Tiebreaker for cursor pagination
});

// 2. get RECENT posts from 'y' days
db.stories.createIndex({
  created_at: -1,        // Primary sort by creation time
  _id: 1                 // Tiebreaker for cursor pagination
});


//COMMENT INDICES//////////////////////////////////////////////////////////////
// 1. get top comments sorted by upvotes on 'x' post
db.comments.createIndex({
  post_id: 1,            // Filter by post
  upvote_count: -1,      // Sort by votes
  _id: 1                 // Tiebreaker for cursor pagination
});

// 2. get RECENT comments on 'x' post
db.comments.createIndex({
  post_id: 1,            // Filter by post
  created_at: -1,        // Sort by time
  _id: 1                 // Tiebreaker for cursor pagination
});
