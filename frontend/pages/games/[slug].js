import Layout from '../../components/layout';
 
export default function Game() {
    return (
        <Layout>
            this is a game
        </Layout>
    )
}
 
export async function getStaticPaths() {
  // Return a list of possible value for id
}
 
export async function getStaticProps({ params }) {
  // Fetch necessary data for the blog post using params.id
}