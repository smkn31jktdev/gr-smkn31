<script lang="ts">
  import { onMount } from 'svelte';

  let currentTime = $state(new Date());

  $effect(() => {
    const interval = setInterval(() => {
      currentTime = new Date();
    }, 1000);
    return () => clearInterval(interval);
  });

  // Derived Clock values
  const clockDigital = $derived.by(() => {
    const hh = String(currentTime.getHours()).padStart(2, '0');
    const mm = String(currentTime.getMinutes()).padStart(2, '0');
    const ss = String(currentTime.getSeconds()).padStart(2, '0');
    return `${hh}:${mm}:${ss}`;
  });

  const clockDateString = $derived.by(() => {
    const options: Intl.DateTimeFormatOptions = { 
      weekday: 'long', 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    };
    return currentTime.toLocaleDateString('id-ID', options);
  });

  // Analog Clock angles
  const clockAngles = $derived.by(() => {
    const h = currentTime.getHours() % 12;
    const m = currentTime.getMinutes();
    const s = currentTime.getSeconds();
    
    return {
      hour: (h * 30) + (m * 0.5),
      minute: (m * 6) + (s * 0.1),
      second: s * 6
    };
  });
</script>

<div class="flex flex-col items-center py-4 w-full">
  <!-- Animated SVG Analog Clock Face -->
  <div class="relative w-44 h-44 mb-6 bg-white rounded-full border border-gray-100 shadow-[inset_0_4px_12px_rgba(0,0,0,0.03),0_8px_24px_rgba(0,0,0,0.05)] flex items-center justify-center">
    <svg class="w-full h-full" viewBox="0 0 100 100">
      <!-- Dial markers -->
      <circle cx="50" cy="50" r="46" fill="none" stroke="var(--color-border)" stroke-width="0.5" />
      
      {#each Array(12) as _, i}
        {@const angle = i * 30}
        {@const rad = (angle * Math.PI) / 180}
        {@const x1 = 50 + 41 * Math.cos(rad)}
        {@const y1 = 50 + 41 * Math.sin(rad)}
        {@const x2 = 50 + 45 * Math.cos(rad)}
        {@const y2 = 50 + 45 * Math.sin(rad)}
        <line x1={x1} y1={y1} x2={x2} y2={y2} stroke={i % 3 === 0 ? 'var(--color-primary)' : 'var(--color-muted)'} stroke-width={i % 3 === 0 ? '1.5' : '0.8'} stroke-linecap="round" />
      {/each}

      <!-- Center dot -->
      <circle cx="50" cy="50" r="2.5" fill="var(--color-primary)" />

      <!-- Hour Hand -->
      <line
        x1="50" y1="50" x2="50" y2="24"
        stroke="var(--color-foreground)" stroke-width="2.5" stroke-linecap="round"
        style="transform: rotate({clockAngles.hour}deg); transform-origin: 50px 50px; transition: transform 0.2s cubic-bezier(0.4, 2.08, 0.55, 1);"
      />
      <!-- Minute Hand -->
      <line
        x1="50" y1="50" x2="50" y2="16"
        stroke="var(--color-muted)" stroke-width="1.8" stroke-linecap="round"
        style="transform: rotate({clockAngles.minute}deg); transform-origin: 50px 50px; transition: transform 0.2s cubic-bezier(0.4, 2.08, 0.55, 1);"
      />
      <!-- Second Hand -->
      <line
        x1="50" y1="55" x2="50" y2="12"
        stroke="var(--color-secondary)" stroke-width="1" stroke-linecap="round"
        style="transform: rotate({clockAngles.second}deg); transform-origin: 50px 50px; transition: transform 0.1s linear;"
      />
      
      <!-- Cap pin -->
      <circle cx="50" cy="50" r="1.2" fill="white" />
    </svg>
  </div>

  <!-- Digital Display -->
  <div class="text-3xl font-black tracking-tight text-gray-800 font-display mb-1.5 bg-gray-50 border border-gray-100 px-5 py-1.5 rounded-2xl shadow-inner">
    {clockDigital}
  </div>
  
  <!-- Date Display -->
  <div class="text-xs font-bold text-teal-655 uppercase tracking-widest text-center max-w-[90%] truncate">
    {clockDateString}
  </div>
</div>
