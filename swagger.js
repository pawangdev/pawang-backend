const swaggerAutogen = require('swagger-autogen')()

const doc = {
    info: {
        title: 'Pawang Rest API',
        version: '3.0.0'
    },
    host: 'api.pawang.studio',
    basePath: '/api',
    schemes: ['https'],
};

const outputFile = './swagger_output.json'
const endpointsFiles = ['./src/routes/index.js']

swaggerAutogen(outputFile, endpointsFiles, doc)
