import gsap from 'gsap'
import constants from '../constants'

const generateMoles = (moles) =>
  moles.map((mole) => ({
    id: mole.id,
    ready: mole.running,
    speed: gsap.utils.random(5, 5.5),
    delay: gsap.utils.random(0.1, 0.2),
    points: constants.MOLE_SCORE,
  }))

export default generateMoles
