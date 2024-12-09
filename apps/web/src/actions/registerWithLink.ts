"use server"

export async function registerWithLink(data: {
  email: string
  password: string
}) {
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
      const { error } = await response.json()
      throw new Error(error || "Registration failed")
    }

    return true
  } catch (error) {
    console.log(error)
    return false
  }
}
