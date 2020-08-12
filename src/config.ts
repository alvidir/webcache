// Configure environment variables
import * as dotenv from 'dotenv';
dotenv.config();

export default {
    ServicePort: process.env.SERVICE_PORT ?? '3001',
    ApiKey: process.env.API_KEY ?? 'None'
}