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

subscriber/<subscriber_email>/newsletter:
  - <newsletter_public_id>: boolean

newsletter/
  - <newsletter_public_id>: boolean

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

## Run All Migrations Down
```makefile
make migration-down-all
```

## Run Migration Down By One
```makefile
make migration-down-by-one
```

### Alternative solution:

#### Whole Newsletter model in Firebase
Move data model in PostgreSQL to the Firebase.

That way subscribers to the newsletter can detail the newsletter by referencing the newsletter by newsletter_id in the Firebase.

**PROS**:
- detail of newsletter is accessed directly from the firebase
- not hiddend behind the backend API

**CONS**:
- overhead of configuring the permission to readonly for subscribers
- listing newsletters for editors would be too much of an overhead -> filtering/pagination of newsletters for an editor would need to be done on the application level and buffered in the memory (potential memory leak if not handled correctly)

## API Docs
- implemented with Swagger UI
- [API Docs](http://165.22.22.96:3000/api/indexlhtml)

## Observability and Health Checks
- [Health Check Readiness Probe](http://165.22.22.96:3000/health/readiness)
- [Health Check Liveness Probe](http://165.22.22.96:3000/health/liveness)
- [Metrics](http://165.22.22.96:3000/metrics)

## TODOs and Missing parts/Enhancements
- Add DNS records for the ip address so it’s reachable via a readable name and not IP.
- Add Cloudflare with TLS protocol with Cloudflare so the communication between client and server is encrypted.
- Add load balancer with nginx and rate limiter so the server is not overloaded.
- Cache docker images in Github actions so it takes less time to run them.
- Add more tests so every reasonable and possible use case is covered with tests.
- Improve monitoring - scrape prometheus metrics and send them into Grafana to observe metrics
- Add ELK stack to gather logs and present them in a user friendly way via Kibana.
- Add sentry instead of ELK stack to cover presenting logs and alerting.
- Enhance Firebase Realtime database ACL rules to specify more granular rules to access the data.
- Replace docker swarm with Kubernetes cluster with helm charts so the deployment and app management is easier and more smooth, better way to manage secrets (use sealed secrets - encrypt only inside the k8s cluster). But questionable because the k8sw cluster needs more resources -> not needed for such a small app.
- For creating subscription verify email addresses before sending the subscription confirmation so the subscription is created only for valid and verified addresses -> could be implemented via webhook, later if the app grows some message broker could be used instead.
- Add domain error mapping to better inform clients about the errors.
