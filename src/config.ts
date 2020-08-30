// Configure environment variables
require('dotenv').config();

export default {
    AppHost: process.env.APP_HOST ?? '127.0.0.1',
    RestServicePort: process.env.REST_SERVICE_PORT ?? '3030',
    ApiKey: process.env.API_KEY ?? 'None',
    ApiTimeout: process.env.API_TIMEOUT ?? '500',
    CacheSize: process.env.CACHE_SIZE ?? '15',
    CacheDeadline: process.env.CACHE_DEADLINE ?? '3600',
}