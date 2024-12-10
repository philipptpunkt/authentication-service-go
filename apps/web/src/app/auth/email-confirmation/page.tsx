export default async function EmailConfirmationPage({
  searchParams,
}: {
  searchParams: { token?: string }
}) {
  const token = searchParams?.token

  if (!token) {
    return (
      <div className="max-w-[500px] mx-auto text-center">
        <h1>Email Confirmation</h1>
        <p className="text-red-500">Invalid or missing confirmation token.</p>
        <a href="/auth/login" className="mt-4 text-blue-500 underline">
          Go to Login
        </a>
      </div>
    )
  }

  const response = await fetch(
    `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/auth/email-confirmation`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ token }),
    }
  )

  if (response.ok) {
    return (
      <div className="max-w-[500px] mx-auto text-center">
        <h1>Email Confirmed</h1>
        <p className="text-green-500">
          Your email has been successfully confirmed!
        </p>
        <a href="/auth/login" className="mt-4 text-blue-500 underline">
          Go to Login
        </a>
      </div>
    )
  }

  const { error } = await response.json()
  return (
    <div className="max-w-[500px] mx-auto text-center">
      <h1>Email Confirmation</h1>
      <p className="text-red-500">
        {error || "Confirmation failed. Please try again."}
      </p>
      <a href="/auth/login" className="mt-4 text-blue-500 underline">
        Go to Login
      </a>
    </div>
  )
}
