import  { GetSignedUrlConfig, Storage } from '@google-cloud/storage'
import express, { Request, Response } from "express"  


const app = express()
const storage = new Storage()

const ttl_mins = 15
const port = process.env.PORT || 8080 


app.use(express.json())
app.post('/', async (req: Request<{}, {}, { bucket: string; filename: string }>, res: Response) => {
    const { bucket, filename } = req.body;
    const url = await generate_signedURLv4(bucket, filename, ttl_mins);
  
    res.json({ url });
});

app.get('/health', (req: Request,res: Response)=>{
    console.log(`GET /health called!`)
    res.send(`Server is healthy, running@PORT:${port}`)
})
  

async function generate_signedURLv4(bucketName: string, fileName: string, ttl_mins: number) : Promise<string>{
    const options : GetSignedUrlConfig = {
        version: 'v4', action: 'read', expires: Date.now() + (ttl_mins * 60 * 1000)
    };
    const [url] = await storage.bucket(bucketName).file(fileName).getSignedUrl(options); 
    return url; 
}



app.listen(port, ()=>{
    console.log(`Server successfully connected @PORT:${port}`)
}).on('error', (err)=>{
    console.log(err);
    throw err;
})