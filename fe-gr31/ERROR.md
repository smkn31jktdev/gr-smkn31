# Fix: `effect_orphan` — Svelte 5 (`$effect.pre` can only be used inside an effect)

## Gejala
- Web blank putih total, gak ada UI yang render
- Console error:
  ```
  Uncaught (in promise) Svelte error: effect_orphan
  `$effect.pre` can only be used inside an effect (e.g. during component initialisation)
  https://svelte.dev/e/effect_orphan
      in <unknown>
      at Root (root.svelte:16:3)
  ```

## Penyebab
Svelte 5 runes (`$effect`, `$effect.pre`) **hanya boleh dipanggil langsung di top-level `<script>` saat component pertama kali di-initialize**. Kalau dipanggil di luar konteks itu, Svelte gak tau effect ini harus "nempel" ke component mana → error `effect_orphan`.

Penyebab umum (cek satu-satu):

| # | Penyebab | Contoh salah |
|---|----------|--------------|
| 1 | Dipanggil di dalam callback function | `onMount(() => { $effect.pre(() => {...}) })` |
| 2 | Dipanggil di dalam function biasa / event handler | `function onClick() { $effect(() => {...}) }` |
| 3 | Dipanggil setelah `await` di script component | `await fetch(...); $effect(() => {...})` |
| 4 | Dipanggil di file `.js`/`.ts` biasa | bukan `.svelte` atau `.svelte.js` |
| 5 | Dipanggil kondisional / dalam loop | `if (x) { $effect(() => {...}) }` |
| 6 | Library/dependency pihak ketiga pakai `$effect` tapi versi Svelte mismatch | cek `package.json` |

## Langkah Debug

1. **Cari semua pemakaian `$effect`** di project:
   ```bash
   grep -rn "\$effect" src/
   ```

2. **Cek tiap hasil**, pastikan posisinya **langsung di top-level `<script>` block**, bukan dibungkus apapun:

   ❌ Salah:
   ```svelte
   <script>
     import { onMount } from 'svelte';

     onMount(() => {
       $effect.pre(() => {
         console.log('ini salah, dibungkus callback');
       });
     });
   </script>
   ```

   ✅ Benar:
   ```svelte
   <script>
     $effect.pre(() => {
       console.log('ini benar, langsung di top-level');
     });
   </script>
   ```

3. **Kalau butuh effect di dalam kondisi/loop**, pindahkan kondisinya ke *dalam* effect, bukan effect-nya yang dikondisikan:

   ❌ Salah:
   ```svelte
   <script>
     if (someCondition) {
       $effect(() => { /* ... */ });
     }
   </script>
   ```

   ✅ Benar:
   ```svelte
   <script>
     $effect(() => {
       if (someCondition) {
         /* ... */
       }
     });
   </script>
   ```

4. **Kalau errornya muncul dari `root.svelte:16:3`** (bukan komponen lo sendiri) → kemungkinan besar berasal dari **dependency/library** yang belum kompatibel penuh dengan Svelte 5. Cek:
   ```bash
   npm ls svelte
   ```
   Pastikan semua dependency yang pakai Svelte runes sudah versi Svelte 5-compatible. Update atau downgrade sesuai kebutuhan.

5. **Kalau pakai SvelteKit**, pastikan `$effect` gak ditaruh di:
   - `+layout.ts` / `+page.ts` (ini file load function, bukan component)
   - `hooks.client.ts` / `hooks.server.ts`
   
   `$effect` cuma valid di file `.svelte` atau `.svelte.js`/`.svelte.ts`.

## Quick Checklist
- [ ] Semua `$effect`/`$effect.pre` ada di top-level script, gak dibungkus function lain
- [ ] Gak ada yang dipanggil setelah `await`
- [ ] Gak ada yang di dalam `if`/`for`/`while` langsung
- [ ] File-nya `.svelte` atau `.svelte.js`, bukan `.js`/`.ts` biasa
- [ ] Versi `svelte` di `package.json` konsisten di semua dependency

## Referensi
- https://svelte.dev/e/effect_orphan