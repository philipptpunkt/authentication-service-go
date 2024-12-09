import { Card } from "@authentication-service-go/ui/Cards"
import { cookies } from "next/headers"
import { redirect } from "next/navigation"

export default function RegisterWithLinkSuccessPage() {
  const token = cookies().get("temp_token")

  if (!token) {
    redirect("/register/link")
  }

  cookies().delete("temp_token")

  return (
    <Card className="max-w-[500px]">
      <h1>Success</h1>
      <p>Thank you for creating an account.</p>
      <p>
        You should have received an Email with a confirmation link. After you
        have verified your email address your account will be ready.
      </p>
    </Card>
  )
}
