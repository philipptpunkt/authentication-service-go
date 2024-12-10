import { Card } from "@authentication-service-go/ui/Cards"
import { redirect } from "next/navigation"
import { ChangePasswordForm } from "../../../components/ChangePasswordForm"

export default function ChangePasswordPage() {
  const handleChangePassword = async (data: {
    currentPassword: string
    newPassword: string
  }) => {
    "use server"
    let success = false

    try {
      const response = await fetch(
        process.env.NEXT_PUBLIC_BACKEND_URL + "/api/v1/auth/change-password",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            current_password: data.currentPassword,
            new_password: data.newPassword,
          }),
        }
      )

      if (!response.ok) {
        const error = await response.text()
        throw new Error(error || "Confirmation failed")
      }

      success = true
    } catch (error) {
      console.log(error)
      success = false
    }

    if (success) {
      redirect("/")
    }
  }

  return (
    <Card className="max-w-[500px]">
      <h1>Change Password</h1>
      <ChangePasswordForm onSubmit={handleChangePassword} />
    </Card>
  )
}
