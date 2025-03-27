<script>
    export let size = 24; // Default size, customizable via props
  </script>
  
  <svg
    width="{size}px"
    height="{Math.round(size * 2/3)}px"
    viewBox="0 0 36 24"
    xmlns="http://www.w3.org/2000/svg"
    aria-label="Flag of Spain"
  >
    <!-- Red top stripe (Spain flag) -->
    <rect width="36" height="8" fill="#AA151B" y="0" />
    
    <!-- Yellow middle stripe -->
    <rect width="36" height="8" fill="#F1BF00" y="8" />
    
    <!-- Red bottom stripe -->
    <rect width="36" height="8" fill="#AA151B" y="16" />
    
    <!-- Optional: Coat of Arms (simplified) -->
    <circle cx="18" cy="12" r="4" fill="#AA151B" />
  </svg>
  
  <style>
    svg {
      display: inline-block;
      vertical-align: middle;
    }
  </style>