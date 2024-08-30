import {Navigate} from "react-router-dom"
import {Dashboard} from "@/components/Login/Login.tsx";
//创建路由
const routes = [
    {
        path: '/',
        element: <Navigate to="/index"/>,
    },
    {
        path: '/login',
        element: <Dashboard/>
    }
]
export default routes;
