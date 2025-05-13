import { Storage } from '@google-cloud/storage';
import express from "express";
const app = express();
const storage = new Storage({
    keyFilename: '/home/dev/.keys/user-o1-gcp.json'
});
const ttl_mins = 15;
const port = 8080;
app.use(express.json());
app.get('/', async (req, res) => {
    let bucket = req.body['bucket'];
    let fileName = req.body['filename'];
    let url = await generate_signedURLv4(bucket, fileName, ttl_mins);
    res.json({
        url
    });
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
    throw err;
});
