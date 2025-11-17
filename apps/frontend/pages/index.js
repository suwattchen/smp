import Head from 'next/head';

export default function Home() {
  return (
    <>
      <Head>
        <title>Sunmart Portal</title>
      </Head>
      <main style={{ padding: '2rem', fontFamily: 'Arial, sans-serif' }}>
        <h1>Sunmart Portal</h1>
        <p>Next.js single-node portal behind Kong.</p>
        <section>
          <h2>Routes</h2>
          <ul>
            <li>
              <strong>Frontend</strong>: served at <code>/</code> via Kong
            </li>
            <li>
              <strong>API</strong>: proxied to core-go at <code>/api</code>
            </li>
          </ul>
        </section>
      </main>
    </>
  );
}
