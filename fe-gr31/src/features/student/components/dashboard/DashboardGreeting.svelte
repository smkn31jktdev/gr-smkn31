<script lang="ts">
  import { currentUser } from '../../../../stores/authStore';

  let name = $derived($currentUser?.nama || 'Siswa');
  let kelas = $derived($currentUser?.kelas || 'Kelas');
  let email = $derived($currentUser?.email || '');

  const todayStr = new Date().toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  });
</script>

<div class="card bg-linear-to-r from-primary to-primary-hover border-none text-white relative overflow-hidden">
  <div class="absolute right-0 top-0 -mt-10 -mr-10 w-40 h-40 bg-white/10 rounded-full blur-xl"></div>
  <div class="absolute left-1/3 bottom-0 -mb-8 w-24 h-24 bg-white/5 rounded-full blur-lg"></div>

  <div class="relative z-10 flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
    <div>
      <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-white/20 text-white backdrop-blur-md mb-3">
        {todayStr}
      </span>
      <h2 class="text-2xl md:text-3xl font-extrabold tracking-tight text-white mb-1">
        Halo, {name}!
      </h2>
      <p class="text-white/80 font-medium text-sm">
        {kelas} | {email}
      </p>
    </div>
    
    <div class="bg-white/15 backdrop-blur-md border border-white/10 p-4 rounded-2xl flex flex-col items-center justify-center min-w-[140px] text-center shadow-xs">
      <span class="text-2xl font-black text-amber-300">7 KAIH</span>
      <span class="text-xs font-semibold text-white/90 mt-1">7 Kebiasaan Baik</span>
    </div>
  </div>
</div>
