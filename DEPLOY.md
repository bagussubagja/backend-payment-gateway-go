# Deployment Guide - Render with Supabase

## 1. Setup Supabase Database

1. Create account at [Supabase](https://supabase.com)
2. Create new project
3. Go to Dashboard > Settings > Database, copy connection string
4. Format: `postgresql://postgres:[password]@db.[project-ref].supabase.co:5432/postgres`

## 2. Deploy to Render

### Setup:
1. Fork repository to GitHub
2. Sign up at [Render](https://render.com)
3. Create new Web Service
4. Connect repository
5. Render will auto-detect from `render.yaml`

### Environment Variables:
```
PORT=10000
DB_HOST=db.your-project-ref.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-supabase-password
DB_NAME=postgres
JWT_SECRET_KEY=your-super-secret-jwt-key
JWT_EXPIRATION_HOURS=24h
MIDTRANS_SERVER_KEY=your-midtrans-server-key
MIDTRANS_CLIENT_KEY=your-midtrans-client-key
MIDTRANS_ENVIRONMENT=sandbox
```

## Tips:
- Ensure Supabase connection string is correct
- Test locally with Supabase first
- Monitor logs for debugging
- Use environment variables, don't hardcode credentials