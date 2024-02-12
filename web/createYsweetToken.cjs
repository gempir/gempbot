import { DocumentManager } from '@y-sweet/sdk';

const manager = new DocumentManager(process.env.YSWEET_URL);

async function main() {
    const clientToken = await manager.getOrCreateDocAndToken(process.env.YSWEET_DOC_ID)

    console.log(JSON.stringify(clientToken));
}

main();
