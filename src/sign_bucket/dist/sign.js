import { Storage } from '@google-cloud/storage';
import express from "express";
let storage;
try {
    storage = new Storage();
}
catch (err) {
    console.log(err);
    throw err;
}
const app = express();
const ttl_mins = 15;
const port = process.env.PORT || 8080;
app.use(express.json());
app.post('/', async (req, res) => {
    const { bucket, filename } = req.body;
    const url = await generate_signedURLv4(bucket, filename, ttl_mins);
    res.json({ url });
});
app.get('/health', (req, res) => {
    console.log(`GET /health called!`);
    res.send(`Server is healthy, running@PORT:${port}`);
});
async function generate_signedURLv4(bucketName, fileName, ttl_mins) {
    const options = {
        version: 'v4', action: 'read', expires: Date.now() + (ttl_mins * 60 * 1000)
    };
    const [url] = await storage.bucket(bucketName).file(fileName).getSignedUrl(options);
    return url;
}
app.listen(port, () => {
    console.log(`Server successfully connected @PORT:${port}`);
}).on('error', (err) => {
    console.log(err);
    throw err;
});
