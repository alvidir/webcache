// Configure environment variables
import * as dotenv from 'dotenv';
dotenv.config();

export default {
    AppHost: process.env.APP_HOST ?? 'localhost',
    RestServicePort: process.env.REST_SERVICE_PORT ?? '3001',
    ProtoServicePort: process.env.PROTO_SERVICE_PORT ?? '3002',
    ApiKey: process.env.API_KEY ?? 'None',
    CacheSize: process.env.CACHE_SIZE ?? '15',
    CacheDeadline: process.env.CACHE_DEADLINE ?? '3600',
}