import Link from "next/link"
import { LoginForm } from "../../../components/LoginForm"
import { redirect } from "next/navigation"
import { Card } from "@authentication-service-go/ui/Cards"

export default function LoginPage() {
  async function handleLogin(data: { email: string; password: string }) {
    "use server"
    let tokens = null

    try {
      const response = await fetch(
        process.env.NEXT_PUBLIC_BACKEND_URL + "/api/v1/auth/login",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(data),
        }
      )

      if (!response.ok) {
        const error = await response.text()
        throw new Error(error || "Login failed")
      }

      const tokenResponse = await response.json()

      console.log(">>>> TOKEN", tokenResponse.access_token)
      console.log(">>>> TOKEN", tokenResponse.refresh_token)

      tokens = tokenResponse
    } catch (error) {
      console.log(error)
      tokens = null
    }

    if (tokens) {
      redirect("/")
    }
  }

  return (
    <Card className="max-w-[500px]">
      <h1>Login</h1>
      <LoginForm onSubmit={handleLogin} />
      <Link
        className="mt-8 mb-2 block px-8 py-2 w-full text-center cursor-pointer text-lg"
        href="/auth/reset-password"
      >
        Forgot password
      </Link>
      <hr />
      <Link
        className="mt-8 block px-8 py-2 w-full text-center cursor-pointer rounded-md border-2 border-black text-lg"
        href="/auth/register/link"
      >
        Create new account
      </Link>
    </Card>
  )
}
