package components

import (
	"strings"

	"skene/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// ASCII logo frames for animation
var logoFrames = []string{
	// Frame 1 - Dots sparse
	`
            ....::::::....
        ..::::::::::::::::::..
      .::::.            .::::.
     ::::                  ::::
    ::::                    ::::
    ::::.                  .::::
     :::::.              .:::::
      .::::::..      ..::::::. 
        .:::::::::::::::::.   
            ..::::::::::..
                  .::::         
                 .::::          
    ....        .::::           
    :::::......:::::            
     .:::::::::::.              
        ........                
`,
	// Frame 2 - Dots denser
	`
            ..:::::::::..
        .::::::::::::::::::::. 
      .:::::::..      ..:::::::.
     ::::::              ::::::
    :::::                  :::::
    ::::::                ::::::
     :::::::.          .:::::::
      .::::::::......::::::::. 
        .:::::::::::::::::::.   
           ...::::::::::::..
                  :::::         
                 :::::          
    .....       :::::           
    ::::::.....:::::            
     .::::::::::::.             
        ..........              
`,
	// Frame 3 - Full density
	`
            ..:::::::::::..
        .:::::::::::::::::::::. 
      .::::::::::.    .:::::::::.
     :::::::              :::::::
    ::::::                  ::::::
    :::::::                :::::::
     :::::::::.        .:::::::::
      .:::::::::::::::::::::::.  
        .::::::::::::::::::::.   
           ..::::::::::::::..
              .:::::::::::       
                 ::::::          
    ......      ::::::           
    :::::::...:::::::            
     .::::::::::::::             
        ............             
`,
}

// Big stylized "GENERATING" text
var GeneratingText = `
 ██████╗ ███████╗███╗   ██╗███████╗██████╗  █████╗ ████████╗██╗███╗   ██╗ ██████╗ 
██╔════╝ ██╔════╝████╗  ██║██╔════╝██╔══██╗██╔══██╗╚══██╔══╝██║████╗  ██║██╔════╝ 
██║  ███╗█████╗  ██╔██╗ ██║█████╗  ██████╔╝███████║   ██║   ██║██╔██╗ ██║██║  ███╗
██║   ██║██╔══╝  ██║╚██╗██║██╔══╝  ██╔══██╗██╔══██║   ██║   ██║██║╚██╗██║██║   ██║
╚██████╔╝███████╗██║ ╚████║███████╗██║  ██║██║  ██║   ██║   ██║██║ ╚████║╚██████╔╝
 ╚═════╝ ╚══════╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝╚═╝  ╚═══╝ ╚═════╝ 
`

// Simplified animated S logo matching the design screenshot
var AnimatedSLogo = []string{
	// Frame showing the dotted S pattern from the screenshot
	`                     ..:::::..                         
              .....::::::::::::-                         
         .......:::::::::::::::.                         
       ........             .::::=.                      
      ........                .::::.                     
      .......                                            
       .......                                           
        ..............:::::::::::.                       
          ..................:::::::=.                    
                             .:::::::.                   
                               :::::::                   
                               .::::::                   
              ...               .:::::.                  
            ......             .:::::                    
            ........       ...:::::                      
              ..................:::                      
                  ..............                         `,
}

// GetLogoFrame returns a logo frame based on time offset
func GetLogoFrame(t float64) string {
	frameIndex := int(t*2) % len(AnimatedSLogo)
	return AnimatedSLogo[frameIndex]
}

// RenderAnimatedLogo renders the animated S logo with subtle animation
func RenderAnimatedLogo(t float64) string {
	logo := AnimatedSLogo[0]

	// Apply subtle color variation based on time
	lines := strings.Split(logo, "\n")
	var result []string

	for i, line := range lines {
		// Create wave effect using time
		brightness := 0.5 + 0.5*sinApprox(t*2+float64(i)*0.3)

		// Choose color based on brightness
		var style lipgloss.Style
		if brightness > 0.7 {
			style = styles.ASCIIAnimated
		} else {
			style = styles.ASCII
		}

		result = append(result, style.Render(line))
	}

	return strings.Join(result, "\n")
}

// Simple sine approximation for animation
func sinApprox(x float64) float64 {
	// Simple Taylor series approximation
	for x > 3.14159 {
		x -= 6.28318
	}
	for x < -3.14159 {
		x += 6.28318
	}

	x2 := x * x
	return x * (1 - x2/6 + x2*x2/120)
}

// RenderGeneratingText renders the big GENERATING text
func RenderGeneratingText() string {
	return styles.Accent.Render(GeneratingText)
}

// StaticLogo for non-animated contexts
var StaticLogo = `
    ███████╗██╗  ██╗███████╗███╗   ██╗███████╗
    ██╔════╝██║ ██╔╝██╔════╝████╗  ██║██╔════╝
    ███████╗█████╔╝ █████╗  ██╔██╗ ██║█████╗  
    ╚════██║██╔═██╗ ██╔══╝  ██║╚██╗██║██╔══╝  
    ███████║██║  ██╗███████╗██║ ╚████║███████╗
    ╚══════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═══╝╚══════╝
`

func RenderStaticLogo() string {
	return styles.Accent.Render(StaticLogo)
}

// Mini loading dots animation
func LoadingDots(t float64) string {
	dots := int(t*3) % 4
	return styles.Muted.Render(strings.Repeat("•", dots) + strings.Repeat(" ", 3-dots))
}

// Progress step indicators (●●●○)
func StepIndicator(current, total int) string {
	var result string
	for i := 0; i < total; i++ {
		if i < current {
			result += styles.Accent.Render("●")
		} else {
			result += styles.Muted.Render("○")
		}
		if i < total-1 {
			result += " "
		}
	}
	return result
}
