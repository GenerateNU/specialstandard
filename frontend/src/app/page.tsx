import { ArrowRight, Users } from 'lucide-react'
import Image from 'next/image'
import Link from 'next/link'

export default function Home() {
  return (
    <div className="font-sans grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
      <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
        <Image
          src="/tss.png"
          alt="The Special Standard logo"
          width={180}
          height={38}
          priority
        />
        <div className="text-center sm:text-left">
          <h1 className="text-4xl font-bold mb-4 tracking-tight text-primary">
            Welcome to The Special Standard!
          </h1>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 w-full max-w-2xl">
          <Link
            href="/students"
            className="group p-6 bg-card rounded-xl shadow-lg hover:shadow-xl transition-all duration-200 border border-default hover:bg-card-hover hover:border-hover"
          >
            <div className="flex items-center justify-between mb-3">
              <Users className="w-10 h-10 text-accent" />
              <ArrowRight className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform group-hover:text-accent" />
            </div>
            <h2 className="text-xl font-semibold text-primary mb-2">
              View Students
            </h2>
            <p className="text-secondary text-sm">
              Browse and manage all student records in the system
            </p>
          </Link>

          <div className="p-6 bg-accent-light rounded-xl border border-default opacity-60 cursor-not-allowed">
            <div className="flex items-center justify-between mb-3">
              <div className="w-10 h-10 bg-card-hover rounded-lg"></div>
              <ArrowRight className="w-5 h-5 text-muted" />
            </div>
            <h2 className="text-xl font-semibold text-primary mb-2">
              More Features
            </h2>
            <p className="text-secondary text-sm">
              Coming soon...
            </p>
          </div>
        </div>

        <ol className="font-mono list-inside list-decimal text-sm/6 text-center sm:text-left text-secondary">
          <li className="mb-2 tracking-[-.01em]">
            Get started by editing
            {' '}
            <code className="bg-accent-light font-mono font-semibold px-1 py-0.5 rounded text-primary">
              src/app/page.tsx
            </code>
            .
          </li>
          <li className="tracking-[-.01em]">
            Save and see your changes instantly.
          </li>
        </ol>

        <div className="flex gap-4 items-center flex-col sm:flex-row">
          <a
            className="rounded-full border border-default transition-colors flex items-center justify-center hover:bg-card-hover hover:border-hover font-medium text-sm sm:text-base h-10 sm:h-12 px-4 sm:px-5 w-full sm:w-auto md:w-[158px] text-primary"
            href="https://nextjs.org/docs?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
            target="_blank"
            rel="noopener noreferrer"
          >
            Next.js docs
          </a>
        </div>
      </main>

      <footer className="row-start-3 flex gap-[24px] flex-wrap items-center justify-center">
        <a
          className="flex items-center gap-2 hover:underline hover:underline-offset-4 text-secondary hover:text-accent"
          href="https://nextjs.org/learn?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Image
            aria-hidden
            src="/file.svg"
            alt="File icon"
            width={16}
            height={16}
          />
          Learn
        </a>
        <a
          className="flex items-center gap-2 hover:underline hover:underline-offset-4 text-secondary hover:text-accent"
          href="https://vercel.com/templates?framework=next.js&utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Image
            aria-hidden
            src="/window.svg"
            alt="Window icon"
            width={16}
            height={16}
          />
          Examples
        </a>
        <a
          className="flex items-center gap-2 hover:underline hover:underline-offset-4 text-secondary hover:text-accent"
          href="https://nextjs.org?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
          target="_blank"
          rel="noopener noreferrer"
        >
          <Image
            aria-hidden
            src="/globe.svg"
            alt="Globe icon"
            width={16}
            height={16}
          />
          Go to nextjs.org →
        </a>
      </footer>
    </div>
  )
}
