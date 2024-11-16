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
      required: ['fingerprint', 'created_at', 'last_login'],
      properties: {
        fingerprint: { bsonType: 'string' },
        created_at: { bsonType: 'date' },
        last_login: { bsonType: 'date' }
      }
    }
  }
});

db.createCollection('posts', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['user_id', 'title', 'body', 'created_at', 'upvote_count', 'comment_count'],
      properties: {
        user_id: { bsonType: 'string' },
        title: { bsonType: 'string' },
        body: { bsonType: 'string' },
        created_at: { bsonType: 'date' },
        upvote_count: { bsonType: 'int' },
        comment_count: { bsonType: 'int' },
        comment_ids: {
          bsonType: 'array',
          items: { bsonType: 'string' }
        }
      }
    }
  }
});

db.createCollection('comments', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['post_id', 'user_id', 'body', 'created_at', 'upvote_count', 'is_deleted'],
      properties: {
        post_id: { bsonType: 'string' },
        user_id: { bsonType: 'string' },
        parent_id: { bsonType: ['string', 'null'] },
        body: { bsonType: 'string' },
        created_at: { bsonType: 'date' },
        upvote_count: { bsonType: 'int' },
        is_deleted: { bsonType: 'bool' }
      }
    }
  }
});

db.createCollection('upvotes', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['parent_id', 'parent_type', 'user_id', 'created_at'],
      properties: {
        parent_id: { bsonType: 'string' },
        parent_type: { bsonType: 'string', enum: ['post', 'comment'] },
        user_id: { bsonType: 'string' },
        created_at: { bsonType: 'date' }
      }
    }
  }
});

// Create indexes
db.users.createIndex({ fingerprint: 1 }, { unique: true });
db.posts.createIndex({ created_at: -1, upvote_count: -1 });
db.posts.createIndex({ user_id: 1 });
db.comments.createIndex({ post_id: 1, created_at: -1 });
db.comments.createIndex({ parent_id: 1 });
db.upvotes.createIndex({ parent_id: 1, user_id: 1 }, { unique: true });
db.upvotes.createIndex({ user_id: 1, created_at: -1 });
