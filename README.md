# Requirements:
## User
- create user (automatically sign in after sign up -> create and return JWT token)
- authenticate user (email, password)
- user authorization (middleware)

## Newsletter
- create newsletter (owned by 1 user, consist of: user_id, name, description?)
- list newsletters (owers can list their newsletter) -> configure ACLs
- allow the client application to connect directly to Firebase and list subscribed
newsletters -> configure ACLs per email to be able to list `subscriber/<subscriber_email>` path

## Subscription
- create subscription to a newsletter with an registered/unregistered email, unique link (<api_url>/<newsletter_id>/<escaped_email_address>)
- send email that confirms subscriptions containing link for unsubscription
- unsubscribe from a newsletter

## Firebase
*DB model*:

subscriber/<subscriber_email>:
  - newsletters:
    - <newsletter_id>: boolean
  - create_at: number

## PostgreSQL
*DB model*:

client
 - id: BIGINT, PK
 - email: TEXT, UNIQUE, NOT NULL
```sql
CREATE TABLE client (
    id BIGSERIAL PRIMARY KEY,
    public_id BIGSERIAL UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

newsletter
  - id: BIGINT, PK
  - client_id: BIGINT, FK(client.id), NOT NULL
  - name: TEXT, NOT NULL
  - description: TEXT, DEFAULT('')
```sql
CREATE TABLE newsletter (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT REFERENCES client(id) NOT NULL,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

newsletter_subscriber
  - email TEXT, NOT NULL
  - newsletter_id, BIGINT, NOT NULL, FK(newsletter.id)
  - subscribed_at, TIMESTAMP
  - PK (email, newsletter_id)
```sql
CREATE TABLE newsletter_subscribers (
  email TEXT NOT NULL,
  newsletter_id BIGINT NOT NULL REFERENCES newsletter(id),
  subscribed_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (email, newsletter_id)
);

CREATE INDEX idx_newsletter_id ON newsletter_subscribers (newsletter_id);
CREATE INDEX idx_email ON newsletter_subscribers (email);
```

### Migrations
## Create New Migration
```makefile
make create-migration name=<migration_name>
```

## Run Migrations Up
```makefile
make migration-up
```

## Run Migrations Down
```makefile
make migration-down
```

### Alternative solution:

#### Newsletter model in Firebase
Move data model in PostgreSQL to the Firebase.

That way subscribers to the newsletter can detail the newsletter by referencing the newsletter by newsletter_id in the Firebase.

**PROS**:
- detail of newsletter is accessed directly from the firebase
- not hiddend behind the backend API

**CONS**:
- overhead of configuring the permission to readonly for subscribers
- listing newsletters for editors would be too much of an overhead -> filtering/pagination of newsletters for an editor would need to be done on the application level and buffered in the memory (potential memory leak if not handled correctly)

