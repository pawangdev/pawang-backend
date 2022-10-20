const swaggerAutogen = require('swagger-autogen')()

const doc = {
    info: {
        title: 'Pawang Rest API',
        version: '3.0.0'
    },
    host: 'localhost:5000',
    basePath: '/api',
    schemes: ['http'],
};

const outputFile = './swagger_output.json'
const endpointsFiles = ['./src/routes/index.js']

swaggerAutogen(outputFile, endpointsFiles, doc)