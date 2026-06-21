<script lang="ts">
  import SiswaLoginForm from '../../../features/auth/components/SiswaLoginForm.svelte';
</script>

<div class="min-h-screen bg-white relative overflow-hidden flex flex-col justify-center items-center p-6 select-none font-sans">
  <div class="absolute inset-0 bg-[radial-gradient(#e2e8f0_1.5px,transparent_1.5px)] [background-size:24px_24px] opacity-40 pointer-events-none"></div>

  <div class="absolute top-0 left-0 w-[550px] h-[550px] rounded-full bg-[#29b6f6]/10 blur-[110px] -translate-x-1/4 -translate-y-1/4 pointer-events-none animate-pulse" style="animation-duration: 9s;"></div>
  <div class="absolute bottom-0 right-0 w-[550px] h-[550px] rounded-full bg-[#0070f3]/5 blur-[130px] translate-x-1/4 translate-y-1/4 pointer-events-none animate-pulse" style="animation-duration: 14s;"></div>

  <a
    href="/login"
    class="absolute top-6 left-6 px-4.5 py-2 bg-white border border-slate-200/80 hover:border-slate-300 rounded-full text-xs font-bold text-slate-500 hover:text-slate-700 hover:bg-slate-50 hover:shadow-sm transition-all flex items-center gap-1.5 cursor-pointer z-20"
  >
    ← Kembali
  </a>

  <div class="w-full max-w-[360px] text-center z-10 flex flex-col items-center">
    <img src="/assets/img/navbar.png" alt="SMK Negeri 31 Jakarta Logo" class="h-8 w-auto object-contain mb-6" />

    <!-- Title and subtitle -->
    <h2 class="text-3xl font-black text-[#0070f3] tracking-tight leading-none mb-2.5">
      Selamat Datang
    </h2>
    <p class="text-xs sm:text-[13px] text-gray-400 font-semibold leading-relaxed mb-8">
      Masuk dan mulai terapkan kebiasaan baik setiap hari!
    </p>

    <SiswaLoginForm />
  </div>

  <footer class="absolute bottom-6 text-center z-10">
    <span class="text-[9px] text-gray-400 font-bold uppercase tracking-widest">
      Gerakan Ramah Anak SMK Negeri 31 Jakarta
    </span>
  </footer>
</div>
