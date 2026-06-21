<script lang="ts">
  import { GraduationCap } from 'lucide-svelte';
  import type { Siswa } from '../../../../../auth/types/auth.types';

  let { students }: { students: Siswa[] } = $props();

  let kelasXCount = $derived(students.filter(s => s.kelas && s.kelas.toUpperCase().startsWith('X ')).length);
  let kelasXiCount = $derived(students.filter(s => s.kelas && s.kelas.toUpperCase().startsWith('XI ')).length);
  let kelasXiiCount = $derived(students.filter(s => s.kelas && s.kelas.toUpperCase().startsWith('XII ')).length);
</script>

<!-- Stat 2: Siswa Kelas X -->
<div class="bg-white p-4 rounded-2xl border border-slate-100/80 shadow-xs flex items-center justify-between hover:border-slate-200 transition-all duration-300">
  <div class="space-y-1 text-left">
    <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest">Siswa Kelas X</span>
    <span class="text-lg font-extrabold text-slate-800">{kelasXCount}</span>
  </div>
  <div class="w-9 h-9 rounded-xl bg-slate-900 flex items-center justify-center text-white border border-slate-950">
    <GraduationCap class="w-4.5 h-4.5" />
  </div>
</div>

<!-- Stat 3: Siswa Kelas XI -->
<div class="bg-white p-4 rounded-2xl border border-slate-100/80 shadow-xs flex items-center justify-between hover:border-slate-200 transition-all duration-300">
  <div class="space-y-1 text-left">
    <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest">Siswa Kelas XI</span>
    <span class="text-lg font-extrabold text-slate-800">{kelasXiCount}</span>
  </div>
  <div class="w-9 h-9 rounded-xl bg-indigo-50/50 flex items-center justify-center text-indigo-650 border border-indigo-100/50">
    <GraduationCap class="w-4.5 h-4.5" />
  </div>
</div>

<!-- Stat 4: Siswa Kelas XII -->
<div class="bg-white p-4 rounded-2xl border border-slate-100/80 shadow-xs flex items-center justify-between hover:border-slate-200 transition-all duration-300">
  <div class="space-y-1 text-left">
    <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest">Siswa Kelas XII</span>
    <span class="text-lg font-extrabold text-slate-800">{kelasXiiCount}</span>
  </div>
  <div class="w-9 h-9 rounded-xl bg-teal-50/50 flex items-center justify-center text-teal-650 border border-teal-100/50">
    <GraduationCap class="w-4.5 h-4.5" />
  </div>
</div>
