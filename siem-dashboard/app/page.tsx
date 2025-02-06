import { redirect } from 'next/navigation'


// const AppComponent = ({ Component, pageProps, currentUser }) => {

const Home = () => {
  console.log("[Home] Redirecting to login page");
  redirect('/login')
}

Home.getInitialProps = async (appContext: any) => {
  console.log("[Home] getInitialProps: appContext", appContext);
  return {}
}

export default Home
