<script lang="ts">
  import { onMount } from 'svelte';
  import { createAduan, listAduanSiswa } from '../../logic/chatLogic';
  import { addToast } from '../../../../stores/uiStore';
  import SubmitButton from '../../../shared/components/SubmitButton.svelte';
  import type { Aduan } from '../../types/student.types';
  import { RotateCw, MessageSquare, MessageCircle } from 'lucide-svelte';

  let rooms = $state<Aduan[]>([]);
  let activeRoom = $state<Aduan | null>(null);
  let newMsgText = $state('');
  let listLoading = $state(false);

  async function loadRooms() {
    listLoading = true;
    const res = await listAduanSiswa();
    rooms = res.items;
    
    // Refresh active room data if open
    if (activeRoom) {
      const refreshed = rooms.find(r => r.id === activeRoom?.id);
      if (refreshed) {
        activeRoom = refreshed;
      }
    }
    listLoading = false;
  }

  onMount(() => {
    loadRooms();
  });

  async function handleCreateTicket(handlers: { resolve: () => void; reject: () => void }) {
    if (!newMsgText.trim()) {
      addToast('Pesan tidak boleh kosong', 'warning');
      handlers.reject();
      return;
    }

    const created = await createAduan(newMsgText);
    if (created) {
      newMsgText = '';
      handlers.resolve();
      await loadRooms();
      activeRoom = created;
    } else {
      handlers.reject();
    }
  }

  // Format date helper
  function formatDate(isoStr: string) {
    if (!isoStr) return '';
    const date = new Date(isoStr);
    return date.toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
</script>

<div class="card p-0 grid grid-cols-1 md:grid-cols-3 min-h-[550px] max-h-[650px] overflow-hidden">
  
  <!-- Left Panel: Chat List -->
  <div class="border-r border-border flex flex-col h-full overflow-hidden">
    <div class="p-4 border-b border-border flex items-center justify-between bg-gray-50/50">
      <h3 class="text-sm font-bold text-foreground">Daftar Konsultasi</h3>
      <button 
        onclick={loadRooms} 
        class="p-1 rounded-lg hover:bg-gray-200 text-muted transition-colors flex items-center justify-center"
        aria-label="Segarkan"
      >
        <RotateCw class="w-4 h-4" />
      </button>
    </div>

    <!-- Rooms List -->
    <div class="flex-1 overflow-y-auto p-3 space-y-2">
      {#if listLoading && rooms.length === 0}
        <p class="text-xs text-muted text-center py-8">Memuat chat...</p>
      {:else if rooms.length === 0}
        <div class="text-center py-12 text-muted px-4">
          <MessageSquare class="w-10 h-10 text-slate-300 mx-auto mb-2" />
          <p class="text-xs font-semibold mt-2">Belum ada riwayat konsultasi.</p>
        </div>
      {:else}
        {#each rooms as room}
          <button
            onclick={() => activeRoom = room}
            class="w-full text-left p-3.5 rounded-xl border transition-all duration-150 flex flex-col gap-1.5 cursor-pointer {activeRoom?.id === room.id ? 'border-primary bg-primary-light/30' : 'border-border bg-white hover:bg-gray-50'}"
          >
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-foreground truncate max-w-[120px]">{room.id}</span>
              <span 
                class="inline-block px-2 py-0.5 rounded-full text-xxs font-bold uppercase"
                class:bg-amber-100={room.status === 'open'}
                class:text-amber-800={room.status === 'open'}
                class:bg-blue-100={room.status === 'in_progress'}
                class:text-blue-800={room.status === 'in_progress'}
                class:bg-gray-100={room.status === 'closed'}
                class:text-gray-800={room.status === 'closed'}
              >
                {room.status === 'in_progress' ? 'diproses' : room.status}
              </span>
            </div>
            
            <p class="text-xs text-muted truncate leading-relaxed">
              {room.messages[0]?.isi || ''}
            </p>
            
            <span class="text-xxs text-muted block text-right mt-1">
              {formatDate(room.createdAt)}
            </span>
          </button>
        {/each}
      {/if}
    </div>
  </div>

  <!-- Right Panel: Message Thread / Consultation details -->
  <div class="col-span-2 flex flex-col h-full overflow-hidden bg-gray-50/30">
    {#if activeRoom}
      <!-- Room Header -->
      <div class="p-4 border-b border-border bg-surface flex items-center justify-between">
        <div>
          <h4 class="text-sm font-bold text-foreground flex items-center gap-2">
            Konsultasi #{activeRoom.id}
            {#if activeRoom.adminNama}
              <span class="text-xxs font-medium text-muted">ditangani oleh {activeRoom.adminNama}</span>
            {/if}
          </h4>
          <p class="text-xxs text-muted mt-0.5">Dibuat pada {formatDate(activeRoom.createdAt)}</p>
        </div>
        <button 
          onclick={() => activeRoom = null} 
          class="text-xs font-bold text-muted hover:text-foreground cursor-pointer md:hidden"
        >
          Kembali
        </button>
      </div>

      <!-- Messages Thread -->
      <div class="flex-1 overflow-y-auto p-4 space-y-3.5">
        {#each activeRoom.messages as msg}
          <div class="flex flex-col max-w-[75%]" class:self-end={msg.role === 'student'} class:ml-auto={msg.role === 'student'}>
            <span class="text-xxs text-muted font-bold mb-1 ml-1" class:text-right={msg.role === 'student'}>
              {msg.role === 'student' ? 'Saya' : (activeRoom.adminNama || 'Konselor')}
            </span>
            <div 
              class="p-3.5 rounded-2xl text-sm leading-relaxed"
              class:bg-primary={msg.role === 'student'}
              class:text-white={msg.role === 'student'}
              class:rounded-tr-none={msg.role === 'student'}
              class:bg-surface={msg.role !== 'student'}
              class:text-foreground={msg.role !== 'student'}
              class:border={msg.role !== 'student'}
              class:border-border={msg.role !== 'student'}
              class:rounded-tl-none={msg.role !== 'student'}
            >
              {msg.isi}
            </div>
            <span class="text-xxs text-muted mt-1 ml-1" class:text-right={msg.role === 'student'}>
              {formatDate(msg.timestamp)}
            </span>
          </div>
        {/each}
      </div>

      <!-- Active Room Footer (Info Box since students only open new tickets) -->
      <div class="p-4 border-t border-border bg-surface text-center">
        <p class="text-xs text-muted font-medium">
          * Untuk mengirim pesan baru atau membalas, silakan buka sesi konsultasi baru di bawah ini.
        </p>
      </div>
    {:else}
      <!-- Empty State / Create New Ticket Screen -->
      <div class="flex-1 flex flex-col items-center justify-center p-6 text-center max-w-md mx-auto h-full">
        <MessageCircle class="w-12 h-12 text-slate-300 mx-auto mb-4" />
        <h4 class="text-base font-bold text-foreground mt-3">Konsultasi Baru</h4>
        <p class="text-xs text-muted mt-1 mb-6">Tuliskan pesan Anda kepada guru BK atau wali kelas secara rahasia. Pertanyaan Anda akan direspons secepat mungkin.</p>
        
        <form class="w-full space-y-4" onsubmit={(e) => e.preventDefault()}>
          <textarea
            bind:value={newMsgText}
            placeholder="Tuliskan keluhan atau hal yang ingin Anda konsultasikan di sini..."
            class="input min-h-[120px] py-3 text-sm"
          ></textarea>
          
          <SubmitButton
            label="Kirim Pesan Konsultasi"
            loadingLabel="Mengirim..."
            className="w-full py-3"
            onclick={handleCreateTicket}
          />
        </form>
      </div>
    {/if}
  </div>
</div>
