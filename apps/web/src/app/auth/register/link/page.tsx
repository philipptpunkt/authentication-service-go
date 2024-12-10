import { Card } from "@authentication-service-go/ui/Cards"
import Link from "next/link"
import { RegisterForm } from "../../../../components/RegisterForm"
import { redirect } from "next/navigation"
import { cookies } from "next/headers"

export default function RegisterWithLinkPage() {
  const handleRegister = async (data: { email: string; password: string }) => {
    "use server"
    let success = false

    try {
      const response = await fetch(
        process.env.NEXT_PUBLIC_BACKEND_URL + "/api/v1/auth/register/link",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(data),
        }
      )

      if (!response.ok) {
        const error = await response.text()
        throw new Error(error || "Registration failed")
      }

      success = true
    } catch (error) {
      console.log(error)
      success = false
    }

    if (success) {
      cookies().set("temp_token", "valid", {
        httpOnly: true,
        secure: true,
        maxAge: 60, // Token expires in 60 seconds
      })

      redirect("/auth/register/link/success")
    }
  }

  return (
    <Card className="max-w-[500px]">
      <h1>Register</h1>
      <RegisterForm onSubmit={handleRegister} />
      <div className="h-16" />
      <hr />
      <Link
        className="mt-8 mb-2 block px-8 py-2 w-full text-center cursor-pointer text-lg"
        href="/login"
      >
        You already have an account?
      </Link>
    </Card>
  )
}
