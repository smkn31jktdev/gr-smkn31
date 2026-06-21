<script lang="ts">
  import DropdownChoice from './DropdownChoice.svelte';

  let {
    options,
    value = $bindable(),
    placeholder,
    onchange,
    minWidth = '130px',
    align = 'left',
    class: className = ''
  }: {
    options: { value: any; label: string }[];
    value: any;
    placeholder: string;
    onchange: () => void;
    minWidth?: string;
    align?: 'left' | 'right';
    class?: string;
  } = $props();
</script>

<div style="min-width: {minWidth};" class="w-full sm:w-auto text-left {className}">
  <DropdownChoice
    {options}
    bind:value
    {onchange}
    {placeholder}
    {align}
  />
</div>
